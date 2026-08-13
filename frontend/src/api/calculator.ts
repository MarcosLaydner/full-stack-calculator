export type Operation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'power'
  | 'square_root'
  | 'percentage'

interface CalculateRequest {
  operation: Operation
  a: number
  b?: number
}

interface CalculateResponse {
  result: number
}

interface ErrorResponse {
  error?: string
}

export async function calculate(request: CalculateRequest): Promise<number> {
  const endpoint = request.operation === 'square_root' ? 'square-root' : request.operation
  const operands = request.b === undefined
    ? { a: request.a }
    : { a: request.a, b: request.b }
  const response = await fetch(`/api/${endpoint}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(operands)
  })

  const payload = await response.json() as CalculateResponse & ErrorResponse
  if (!response.ok) {
    throw new Error(payload.error ?? 'Unable to complete the calculation')
  }

  return payload.result
}
