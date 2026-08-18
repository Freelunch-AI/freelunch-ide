# Lunch Football | On-demand football analytics

Today's football analytics platforms answer questions someone decided to measure years ago, with large databases of pre-defined metrics. We answer the questions that are actually in a coach's or scout's head. Ask any specific and complex thing—from tactical behavior to player decision-making—and our AI research agent automatically watches the relevant match footage, builds the analysis, quantifies its findings, and returns an explainable report with linked evidence video clips.

Note: can get free match footage on footballia

Examples of questions:

- *"Is this center back winning aerial duels because of positioning, timing, or simply physical dominance?"*
- *"Show me every time our high press failed because the left winger didn't trigger at the right moment."*
- *"Which fullback is better at delaying counterattacks rather than attempting tackles?"*
- *"Does this midfielder scan before receiving under pressure more consistently than other players in the league?"*
- *"Find wingers similar to this player, but who are better at creating separation before receiving the ball."*
- *"Against top-six teams, does our defensive line become too narrow when defending switches of play?"*
- *"Has this striker improved his movement to attack the near post over the last two seasons?"*

For the example:

> **"Find defenders that intentionally delay rather than tackle."**
> 

Current systems cannot respond this directly, but can provide metrics such as:

| Metric | What it measures |
| --- | --- |
| Tackles attempted | Number of tackles |
| Tackle success % | Successful tackles / attempted tackles |
| Interceptions | Reading passing lanes |
| Blocks | Blocked shots/crosses |
| Clearances | Balls cleared |
| Defensive duels | Number of ground duels |
| Defensive duel success % | Win rate in defensive duels |
| Fouls committed | Discipline |
| Dribblers tackled | Success vs dribblers |
| Times dribbled past | Beaten by an attacker |
| Pressures | Number of pressure actions |
| Pressure success | Possession won after pressure |
| Challenge intensity | How often player engages |
| xT prevented / xG prevented | Some providers estimate defensive value |

These are available (to varying degrees) from providers like StatsBomb, Opta, Wyscout, and Hudl.

---

## Why these don't answer the coach's question

Imagine two fullbacks.

### Player A

```
Tackles: 3.2
Success: 72%
```

### Player B

```
Tackles: 1.0
Success: 95%
```

Which one is better at **delaying attacks?**

Impossible to know.

Maybe Player B

- forces the attacker backwards
- waits for teammates
- prevents dangerous passes
- never needs to tackle

The statistics actually make him look *less active*.

---

## What the coach actually means

When a coach says

> "He delays attackers well."
> 

They're usually referring to things like

- staying goal side
- preventing penetration
- forcing attacker wide
- slowing transition
- waiting for defensive support
- avoiding diving into tackles
- reducing shot probability

Almost none of those are explicit events.

They're **behaviors**.

---

## Another example

Suppose a coach asks

> "Which CBs are patient defenders?"
> 

Current metrics

```
Tackle success
Interceptions
Fouls
Duels won
```

But patience might mean

- waits 2–3 seconds before engaging
- shepherds striker away
- lets passing lane close naturally
- only tackles when probability of success exceeds a threshold

Those aren't tracked.

---

## The bigger distinction

I think the boundary is this:

### Existing systems answer

> **What happened?**
> 

Examples

- 8 tackles
- 5 interceptions
- 78% aerial duel success
- 11 clearances
- 23 pressures

---

### Coaches often ask

> **Why did it happen?**
> 

Examples

- Did he overcommit?
- Did he delay correctly?
- Did he bait the pass?
- Did he anticipate well?
- Was his body orientation correct?
- Did he choose the right defensive action?
- Was he protecting space instead of chasing the ball?

Those are much harder because they require interpreting intent or tactical context.
