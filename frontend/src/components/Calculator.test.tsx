import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Calculator } from './Calculator'

describe('Calculator', () => {
  afterEach(() => {
    cleanup()
    vi.restoreAllMocks()
  })

  it('validates required operands', async () => {
    render(<Calculator />)
    await userEvent.click(screen.getByRole('button', { name: 'Calculate' }))
    expect(screen.getByRole('alert')).toHaveTextContent('Enter a valid first number')
  })

  it('sends a calculation and displays the result', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(JSON.stringify({ result: 42 }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' }
      })
    )
    render(<Calculator />)

    await userEvent.type(screen.getByLabelText('First number'), '35')
    await userEvent.type(screen.getByLabelText('Second number'), '7')
    await userEvent.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(await screen.findByText('42')).toBeInTheDocument()
    expect(fetch).toHaveBeenCalledWith('/api/add', expect.objectContaining({
      method: 'POST',
      body: JSON.stringify({ a: 35, b: 7 })
    }))
  })

  it('shows API errors', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(JSON.stringify({ error: 'division by zero' }), {
        status: 422,
        headers: { 'Content-Type': 'application/json' }
      })
    )
    render(<Calculator />)

    await userEvent.click(screen.getByRole('button', { name: 'Divide' }))
    await userEvent.type(screen.getByLabelText('First number'), '5')
    await userEvent.type(screen.getByLabelText('Second number'), '0')
    await userEvent.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('division by zero')
  })
})
