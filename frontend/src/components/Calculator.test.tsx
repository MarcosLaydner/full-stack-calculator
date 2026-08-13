import { act, cleanup, render, screen } from '@testing-library/react'
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

  it('keeps the calculator cleared when an in-flight request resolves', async () => {
    let resolveRequest!: (response: Response) => void
    vi.spyOn(globalThis, 'fetch').mockImplementation(() => new Promise((resolve) => {
      resolveRequest = resolve
    }))
    render(<Calculator />)

    await userEvent.type(screen.getByLabelText('First number'), '35')
    await userEvent.type(screen.getByLabelText('Second number'), '7')
    await userEvent.click(screen.getByRole('button', { name: 'Calculate' }))

    const requestOptions = vi.mocked(fetch).mock.calls[0][1]
    expect(requestOptions?.signal?.aborted).toBe(false)
    await userEvent.click(screen.getByRole('button', { name: 'Clear' }))

    expect(requestOptions?.signal?.aborted).toBe(true)
    expect(screen.getByLabelText('First number')).toHaveValue(null)
    expect(screen.getByLabelText('Second number')).toHaveValue(null)
    expect(screen.getByText('0')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Calculate' })).toBeEnabled()

    await act(async () => {
      resolveRequest(new Response(JSON.stringify({ result: 42 }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' }
      }))
    })

    expect(screen.getByText('0')).toBeInTheDocument()
    expect(screen.queryByText('42')).not.toBeInTheDocument()
  })

  it('ignores an in-flight result after the operation changes', async () => {
    let resolveRequest!: (response: Response) => void
    vi.spyOn(globalThis, 'fetch').mockImplementation(() => new Promise((resolve) => {
      resolveRequest = resolve
    }))
    render(<Calculator />)

    await userEvent.type(screen.getByLabelText('First number'), '35')
    await userEvent.type(screen.getByLabelText('Second number'), '7')
    await userEvent.click(screen.getByRole('button', { name: 'Calculate' }))
    const requestOptions = vi.mocked(fetch).mock.calls[0][1]
    await userEvent.click(screen.getByRole('button', { name: 'Multiply' }))

    expect(requestOptions?.signal?.aborted).toBe(true)
    expect(screen.getByRole('button', { name: 'Multiply' })).toHaveAttribute('aria-pressed', 'true')

    await act(async () => {
      resolveRequest(resultResponse(42))
    })

    expect(screen.getByText('0')).toBeInTheDocument()
    expect(screen.queryByText('42')).not.toBeInTheDocument()
  })

  it('ignores an in-flight result after an operand changes', async () => {
    let resolveRequest!: (response: Response) => void
    vi.spyOn(globalThis, 'fetch').mockImplementation(() => new Promise((resolve) => {
      resolveRequest = resolve
    }))
    render(<Calculator />)

    await userEvent.type(screen.getByLabelText('First number'), '35')
    await userEvent.type(screen.getByLabelText('Second number'), '7')
    await userEvent.click(screen.getByRole('button', { name: 'Calculate' }))
    const requestOptions = vi.mocked(fetch).mock.calls[0][1]
    await userEvent.type(screen.getByLabelText('First number'), '0')

    expect(requestOptions?.signal?.aborted).toBe(true)
    expect(screen.getByLabelText('First number')).toHaveValue(350)

    await act(async () => {
      resolveRequest(resultResponse(42))
    })

    expect(screen.getByText('0')).toBeInTheDocument()
    expect(screen.queryByText('42')).not.toBeInTheDocument()
  })
})

function resultResponse(result: number): Response {
  return new Response(JSON.stringify({ result }), {
    status: 200,
    headers: { 'Content-Type': 'application/json' }
  })
}
