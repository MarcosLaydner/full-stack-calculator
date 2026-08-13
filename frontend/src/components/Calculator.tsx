import { useEffect, useRef, useState } from 'react'
import { calculate, type Operation } from '../api/calculator'

const operations: Array<{ value: Operation; symbol: string; label: string }> = [
  { value: 'add', symbol: '+', label: 'Add' },
  { value: 'subtract', symbol: '\u2212', label: 'Subtract' },
  { value: 'multiply', symbol: '\u00d7', label: 'Multiply' },
  { value: 'divide', symbol: '\u00f7', label: 'Divide' },
  { value: 'power', symbol: 'x\u02b8', label: 'Power' },
  { value: 'percentage', symbol: '%', label: 'Percentage' },
  { value: 'square_root', symbol: '\u221a', label: 'Square root' }
]

function formatResult(value: number): string {
  return Number.isInteger(value) ? value.toString() : Number(value.toPrecision(12)).toString()
}

export function Calculator() {
  const [operation, setOperation] = useState<Operation>('add')
  const [first, setFirst] = useState('')
  const [second, setSecond] = useState('')
  const [result, setResult] = useState<string>('0')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const activeRequest = useRef<AbortController | null>(null)
  const requestID = useRef(0)
  const isUnary = operation === 'square_root'

  useEffect(() => () => activeRequest.current?.abort(), [])

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError('')

    const a = Number(first)
    const b = Number(second)
    if (first.trim() === '' || !Number.isFinite(a)) {
      setError('Enter a valid first number')
      return
    }
    if (!isUnary && (second.trim() === '' || !Number.isFinite(b))) {
      setError('Enter a valid second number')
      return
    }

    activeRequest.current?.abort()
    const controller = new AbortController()
    const currentRequestID = ++requestID.current
    activeRequest.current = controller
    setIsLoading(true)
    try {
      const value = await calculate(
        { operation, a, ...(!isUnary && { b }) },
        controller.signal
      )
      if (currentRequestID === requestID.current) {
        setResult(formatResult(value))
      }
    } catch (caughtError) {
      if (currentRequestID === requestID.current && !controller.signal.aborted) {
        setError(caughtError instanceof Error ? caughtError.message : 'Something went wrong')
      }
    } finally {
      if (currentRequestID === requestID.current) {
        activeRequest.current = null
        setIsLoading(false)
      }
    }
  }

  function clear() {
    requestID.current++
    activeRequest.current?.abort()
    activeRequest.current = null
    setFirst('')
    setSecond('')
    setResult('0')
    setError('')
    setIsLoading(false)
  }

  return (
    <section className="calculator" aria-label="Calculator">
      <div className="display" aria-live="polite">
        <span className="display-label">Result</span>
        <output>{result}</output>
      </div>

      <form onSubmit={handleSubmit} noValidate>
        <fieldset>
          <legend>Choose an operation</legend>
          <div className="operation-grid">
            {operations.map((item) => (
              <button
                className={operation === item.value ? 'operation active' : 'operation'}
                key={item.value}
                type="button"
                aria-label={item.label}
                aria-pressed={operation === item.value}
                onClick={() => {
                  setOperation(item.value)
                  setError('')
                }}
              >
                {item.symbol}
              </button>
            ))}
          </div>
        </fieldset>

        <div className={isUnary ? 'inputs unary' : 'inputs'}>
          <label>
            <span>{isUnary ? 'Number' : 'First number'}</span>
            <input
              inputMode="decimal"
              type="number"
              step="any"
              value={first}
              onChange={(event) => setFirst(event.target.value)}
              placeholder="0"
            />
          </label>
          {!isUnary && (
            <label>
              <span>Second number</span>
              <input
                inputMode="decimal"
                type="number"
                step="any"
                value={second}
                onChange={(event) => setSecond(event.target.value)}
                placeholder="0"
              />
            </label>
          )}
        </div>

        <div className="feedback" role="alert">{error}</div>
        <div className="actions">
          <button className="clear" type="button" onClick={clear}>Clear</button>
          <button className="calculate" type="submit" disabled={isLoading}>
            {isLoading ? 'Calculating...' : 'Calculate'}
          </button>
        </div>
      </form>
    </section>
  )
}
