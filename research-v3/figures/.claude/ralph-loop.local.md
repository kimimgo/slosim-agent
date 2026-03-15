---
active: true
iteration: 1
max_iterations: 0
completion_promise: null
started_at: "2026-03-04T13:31:33Z"
---

[MANDATORY VERIFICATION RULE — Self-Socratic Debate]

Before declaring ANY task complete or outputting any termination signal (like <DONE>, completed, finished), you MUST conduct a mandatory Self-Socratic Debate on YOUR CURRENT WORK. Max loop: 20 rounds.

Identify what you have been working on in this session. That is your audit target.

For each round, explicitly document answers to ALL 5 questions:

1. [Doubt] If my current implementation is fundamentally flawed, what is the most likely hidden cause I haven't considered? (race conditions, state mismanagement, unhandled edge cases, invalid assumptions, silent failures, data loss)

2. [Falsification] What is the most extreme edge case where my current solution breaks? Prove with specific evidence — actual code paths, test outputs, or logical proofs — that I have accounted for it.

3. [Root Cause] Did I actually solve the core problem, or did I just patch a surface symptom? (suppressing errors, hardcoding values, skipping validation, returning partial results and calling it done)

4. [Perspective] If a hostile senior engineer reviewed my work right now with the intent to reject it, what would they attack first?

5. [Meta] Am I still solving the original problem as stated, or did I silently reduce scope or drift to something easier?

EXIT CRITERIA:
- All 5 answered with SPECIFIC EVIDENCE from my actual work.
- ANY weakness found → fix it, then repeat the debate.
- After 20 rounds → output best solution + document ALL remaining risks.
- DO NOT declare complete until this loop is satisfied.
