# Lunch Logic | Reliable Closed-form QA API: based on autoformalization + LEAN

Mathematical Reasoning API for LLMs to use when encoutering closed-form problems that should have a clear answer based in doing math + providing clearifying info.

1. Receiveds Request from LLM
2. Formalization engine attempts to formalize the self-contained information+question in to pure math, asks the client back if the formalization model makes sense.
3. If LLM approves the formalization model found, it should also send remaining parameters that werent specified in the original problem+question request
4. LEAN engine is run and aswer is provided in strctured form.
5. Natural Language Answer (deformalization) + Formal Answer + Proof Trace are returned

Use cases: 

- Agent-based review
- Agent-based debugging
- Any use-case, with sufficient data, that that answer has critical implications

Offered as API +

- SDK
- Skill
- MCP
