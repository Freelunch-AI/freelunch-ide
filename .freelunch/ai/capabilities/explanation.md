# Explanation And Developer Handoff

Use this capability when the developer asks to understand the change or when a risky decision needs an explicit knowledge handoff.

1. Explain the change from the issue and acceptance criteria down to the affected components and data flow.
2. Show important tradeoffs, failure paths, security boundaries, and how tests demonstrate behavior.
3. Ground the explanation in actual files and observed evidence.
4. Ask targeted questions only when the developer requests a knowledge check or an approval depends on understanding a risk.
5. Correct misunderstandings directly and link back to the relevant code or decision.

Do not turn a developer quiz into a mandatory PR gate. A concise, accurate handoff is enough for straightforward changes.
