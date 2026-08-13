import { Calculator } from './components/Calculator'

export function App() {
  return (
    <main className="page-shell">
      <header className="app-header">
        <span className="brand-mark" aria-hidden="true">=</span>
        <div>
          <p className="eyebrow">Full-stack calculator</p>
          <h1>Make it count.</h1>
        </div>
      </header>
      <Calculator />
      <footer>Calculated by a Go microservice</footer>
    </main>
  )
}
