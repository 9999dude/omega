# What an SRE manager does in a FAANG-style company

An SRE manager is responsible for three interconnected systems:

1. **The production system** — availability, latency, capacity, scalability and recovery.
2. **The engineering system** — monitoring, deployment safety, incident response, automation and architecture.
3. **The human system** — engineers, on-call health, skills, ownership, prioritisation and collaboration.

Their real job is not:

> “Make sure there are no outages.”

It is:

> **Create an engineering organisation that can safely operate, change and scale important systems without depending on constant human intervention or heroics.**

Google describes SRE as applying software engineering to operational problems and deliberately limits operational work so that SREs retain enough time for engineering. Its published model caps aggregate operational work at 50%, leaving at least half of SRE time for engineering projects. ([Google SRE][1])

---

## 1. “SRE manager at FAANG” is not one uniform role

Different companies organise reliability differently:

* **Google:** Dedicated SRE teams may support important services and use SLOs, error budgets, production-readiness reviews and engineering-driven operations.
* **Meta:** The equivalent role is often associated with Production Engineering. Production engineers combine software and systems engineering and focus on reliability, scalability, performance and security while partnering with product engineering teams. ([Engineering at Meta][2])
* **Amazon:** Individual service teams generally have strong operational ownership. Reliability management may sit within a service engineering organisation rather than a separate SRE organisation.
* **Netflix:** Product teams generally operate what they build, while central reliability teams provide common tooling and coordinate major incidents affecting customer experience. ([Netflix TechBlog][3])

Therefore, the title may be:

* SRE Manager
* Production Engineering Manager
* Infrastructure Engineering Manager
* Reliability Engineering Manager
* Platform Engineering Manager
* Engineering Manager for a service-owning team

But the underlying responsibility is similar: **manage operational risk using engineering rather than manual effort.**

---

# 2. The exact role of an SRE manager

## 2.1 Own reliability outcomes

The manager ensures that the team and its partner teams can answer:

* What user journeys are critical?
* How is reliability measured?
* What level of reliability is required?
* How much unreliability is acceptable?
* Which failures are currently most dangerous?
* What should the organisation invest in next?

The team might manage:

* Availability
* Request success rate
* Latency
* Data freshness
* Correctness
* Durability
* Capacity
* Deployment safety
* Recovery time
* Operational workload

Google defines traditional SRE responsibility around availability, latency, performance, efficiency, change management, monitoring, emergency response and capacity planning. More recent product-oriented approaches recommend connecting these measurements to critical user functionality rather than looking only at individual services. ([Google SRE][4])

### Average manager

> “Our Kubernetes API servers had 99.99% uptime.”

### Good manager

> “Could application teams deploy, scale and recover workloads successfully? Which critical user journeys were affected when the API server was technically available but admission webhooks were timing out?”

A component can be “green” while users are unable to complete their work.

---

## 2.2 Establish SLOs and error-budget policies

The manager does not personally write every SLI query, but ensures that:

* SLOs represent user experience.
* Product and engineering agree on the target.
* The measurement is trustworthy.
* Error-budget policy is agreed before an incident.
* Reliability data influences prioritisation.

For example:

```text
Critical journey:
Deploy a workload to a production cluster.

SLI:
Percentage of valid deployment requests reaching a ready state
within 10 minutes.

SLO:
99.9% over 28 days.

Error-budget policy:
If 50% of the monthly budget is consumed in seven days:
- pause risky control-plane rollouts;
- assign owners to the dominant failure mode;
- require additional canary validation;
- review whether capacity or dependency failures are involved.
```

An error budget is the permitted unreliability implied by the SLO. For a 99.9% target, the error budget is 0.1%. Google’s published guidance uses error budgets to make the reliability-versus-release-velocity decision more objective. ([Google SRE][5])

The SRE manager’s role is not to make reliability absolute. It is to make the trade-off explicit.

---

## 2.3 Build a sustainable on-call system

The manager is accountable for the health of the on-call system, including:

* Rotation size
* Primary and secondary coverage
* Handover quality
* Escalation paths
* Alert quality
* Runbooks
* Training
* Shadow rotations
* Incident authority
* Recovery time after severe incidents
* Psychological safety
* Fair distribution of operational burden

They should know:

* Pages per engineer per shift
* After-hours pages
* Percentage of actionable pages
* Repeat-page rate
* Time to acknowledge
* Time to mitigate
* Escalation frequency
* Services without trained responders
* Operational work by engineer
* Burnout and attrition signals

The purpose of on-call is not to prove that engineers are committed. It is to provide a safe human fallback when automation cannot handle an abnormal condition.

### Average manager

> “The rotation is working. Someone always responds.”

### Good manager

> “Why does a human need to respond to this condition? Can we prevent it, automatically mitigate it, reduce its blast radius or convert it from a page into a ticket?”

Google classifies operational load into pages, tickets and ongoing operational activities, and stresses that these interrupts need explicit management rather than being treated as invisible work. ([Google SRE][6])

---

## 2.4 Turn operational pain into engineering work

This is one of the defining responsibilities.

An average operations organisation absorbs increasing manual work by adding people.

A good SRE organisation asks:

> “What software, architecture or ownership change would stop this work growing linearly?”

Google defines toil as work that is typically manual, repetitive, tactical, automatable, without enduring value and likely to scale with service growth. ([Google SRE][7])

Examples:

| Operational symptom                  | Weak response               | Engineering response                                             |
| ------------------------------------ | --------------------------- | ---------------------------------------------------------------- |
| Frequent disk cleanup                | Expand the runbook          | Fix retention, quotas and lifecycle automation                   |
| Repeated certificate expiry          | Add calendar reminders      | Build automated issuance, renewal and expiry validation          |
| Manual cluster upgrades              | Assign upgrade coordinators | Create progressive upgrade orchestration with health gates       |
| Frequent pod OOM alerts              | Restart the pod             | Fix sizing, memory behaviour, autoscaling and admission controls |
| Teams request access manually        | Hire more support engineers | Build policy-based self-service workflows                        |
| Deployment rollback requires experts | Document the commands       | Add automatic rollback and safe rollout controls                 |

The SRE manager must protect time for these investments. Otherwise, urgent work continuously displaces important work.

---

## 2.5 Set technical direction without becoming the technical bottleneck

An SRE manager needs enough technical depth to evaluate:

* Failure modes
* Capacity models
* Distributed-system risks
* Data-loss scenarios
* Rollout safety
* Dependency behaviour
* Regional isolation
* Recovery mechanisms
* Observability design
* Control-plane versus data-plane impact
* Security and access risks

But they should not personally make every design decision.

A healthy division is:

### SRE manager

Owns:

* Why the work matters
* Priority
* Staffing
* Risk acceptance
* Cross-team alignment
* Delivery mechanisms
* Accountability
* Team health
* Escalation

### Tech lead or staff engineer

Owns:

* Technical strategy
* Architecture
* Design quality
* Engineering standards
* Technical decomposition
* Mentorship
* Complex implementation decisions

### Incident commander

Owns:

* The response to a specific active incident
* Mitigation priority
* Coordination
* Decision log
* Communication cadence

One person may temporarily perform multiple roles, but a manager who must personally lead every design and every incident has usually failed to build sufficient organisational capacity.

---

## 2.6 Manage people

The SRE manager is still an engineering manager.

They are responsible for:

* Hiring
* Performance management
* Career development
* Promotions
* Coaching
* Delegation
* Succession planning
* Team composition
* Retention
* Conflict resolution
* Creating technical leadership opportunities

SRE adds several complications:

* Operational work is often less visible than feature work.
* Engineers may receive recognition for heroic incidents instead of prevention.
* On-call burden may be uneven.
* Interrupts can destroy project progress.
* Some engineers become indispensable operational experts.
* Reliability work often requires influence across organisational boundaries.

A good manager corrects these incentive problems.

For example, they recognise an engineer who eliminated 300 annual pages, not only the engineer who handled a dramatic outage at 3 a.m.

---

## 2.7 Manage the relationship with product and development teams

SRE cannot succeed as a team that merely receives broken services.

The manager establishes an engagement model covering:

* What SRE owns
* What the development team owns
* Support boundaries
* Production-readiness requirements
* Escalation paths
* On-call expectations
* Reliability targets
* Launch responsibilities
* Exit criteria
* What happens when operational load becomes excessive

Production-readiness reviews often cover design, monitoring, capacity, failure handling, operational documentation and training before an SRE team assumes production responsibility. ([Google SRE][8])

A weak manager says:

> “SRE owns production.”

A strong manager says:

> “The product team owns the service. SRE provides reliability expertise, shared engineering and specific operational ownership under an explicit agreement.”

Otherwise, SRE becomes the organisation’s dumping ground.

---

# 3. What their day-to-day work looks like

There is no single daily schedule. The work is usually divided among:

* Production and operational health
* People management
* Engineering execution
* Cross-team alignment
* Technical and risk reviews
* Incident handling
* Organisational planning

## Example of a normal day

### 8:30–9:00: Review service health

The manager does not normally stare at dashboards continuously.

They review exceptions and trends:

* SLO burn
* Major alerts
* Open incidents
* Failed changes
* Capacity risks
* Dependency degradation
* On-call summary
* Significant customer impact

Modern operational monitoring should rely on automated alarms rather than requiring people to continuously watch dashboards. Dashboards are primarily for gaining context, investigation and operational review. ([Amazon Web Services, Inc.][9])

The manager asks:

* Was there user harm?
* Is this isolated or systemic?
* Is mitigation complete?
* Does anyone need help?
* Is there a repeated pattern?
* Does the roadmap need to change?

They should not immediately enter every technical discussion.

---

### 9:00–9:30: Team or operational stand-up

Topics might include:

* Current production concerns
* On-call handover
* Blocked engineering work
* Risky upcoming changes
* Capacity events
* Cross-team dependencies
* Projects needing decisions

A good stand-up is not a status recital.

It identifies:

* Decisions
* Risks
* Dependencies
* Ownership gaps
* Places where management intervention is needed

---

### 9:30–11:00: One-to-ones

A manager may discuss:

* Current workload
* Career goals
* On-call experience
* Project ownership
* Technical growth
* Collaboration problems
* Motivation
* Feedback
* Performance expectations

A strong SRE manager also asks:

* Are you being interrupted too often?
* Is operational knowledge concentrated in one person?
* Are you doing work below your level?
* Do you have enough uninterrupted engineering time?
* Is on-call affecting your health?
* What should the team stop doing?
* Which recurring problem should we engineer away?

---

### 11:00–12:00: Technical or production-readiness review

Example questions for a new service:

* What is the critical user journey?
* What is the expected load?
* What happens at 2× or 10× load?
* What is the load-shedding behaviour?
* Which dependencies are hard requirements?
* What is the blast radius of a bad release?
* Can the service roll back?
* Is rollback actually safe after schema changes?
* What happens if one region fails?
* How is data restored?
* Which alerts require immediate human action?
* Has the failure mode been tested?

Safe rollout mechanisms can include canaries, one-box deployments, traffic splitting, rolling deployment and blue-green deployment. ([AWS Documentation][10])

The manager may not design every mechanism, but ensures that the right questions are answered and owners are assigned.

---

### 13:00–14:00: Product or partner-team planning

The manager represents operational reality during product planning.

For example:

```text
Product proposal:
Launch a new workload scheduling feature next quarter.

SRE concerns:
- Scheduler control plane already has high tail latency.
- Regional failover has not been tested at projected scale.
- Two key dependencies have no overload protection.
- Current on-call receives 12 actionable pages per week.

Decision:
Proceed with limited-region launch after:
1. Capacity testing
2. Dependency timeout changes
3. Load shedding
4. Canary metrics
5. Rollback automation
6. Game day
```

The manager is not simply saying “no.” They are converting reliability risk into a sequence of executable conditions.

---

### 14:00–15:00: Project and roadmap review

The manager reviews:

* Outcome achieved, not just tasks completed
* Milestones
* Risk
* Scope
* Staffing
* Dependency delays
* Whether a project still has sufficient value
* Whether production events have changed its priority

For example:

Bad progress report:

> “The automated remediation project is 70% complete.”

Better progress report:

> “The remediation currently handles two of the five dominant failure modes and has removed 18% of last quarter’s pages. Database saturation remains the largest unresolved source.”

---

### 15:00–16:00: Incident or postmortem review

Questions should include:

* What was the user impact?
* Why was detection early or late?
* What made the failure possible?
* What increased the blast radius?
* What slowed mitigation?
* Which assumptions were wrong?
* Which safeguards failed?
* Did organisational structure contribute?
* Are corrective actions proportional to risk?
* Who owns them?
* Which action prevents recurrence rather than merely documenting it?

Google’s incident-management model separates responsibilities such as incident command, operations, planning and communication so responders can coordinate rather than everyone independently debugging. ([Google SRE][11])

---

### 16:00–17:30: Focused management work

This may include:

* Writing a reliability strategy
* Preparing headcount justification
* Performance feedback
* Hiring
* Reviewing an organisational proposal
* Escalating dependency risk
* Drafting quarterly objectives
* Resolving ownership disputes
* Analysing operational metrics
* Planning a game day
* Reviewing an incident trend

A manager’s most valuable work is often written communication and decision-making rather than attending more meetings.

---

# 4. What happens during a major incident

A bad SRE manager becomes the loudest person in the incident channel.

A good SRE manager makes the response system work.

## Bad incident behaviour

The manager:

* Interrupts responders for status every five minutes
* Suggests random fixes
* Changes priorities without understanding the state
* Asks who caused the incident
* Pulls too many people into the call
* Communicates unverified theories to executives
* Keeps engineers working after they are exhausted
* Tries to remain incident commander, technical lead and executive communicator simultaneously

This increases cognitive load.

## Good incident behaviour

The manager confirms:

* There is a named incident commander.
* Someone owns mitigation.
* Someone maintains the incident state document.
* Someone owns communication.
* Customer impact is understood.
* Responders have authority to act.
* Executives are not distracting technical responders.
* Handoffs occur if the incident is prolonged.
* Risky actions have clear decision owners.
* The team is focused on mitigation before root-cause perfection.

The manager may say:

> “Priya is incident commander. Daniel owns database mitigation. Mei owns customer-impact measurement. I will handle leadership communication. Technical responders should route decisions through Priya.”

Netflix describes a similar incident-management function involving coordination, decision-making, incident recording, technical investigation and communication across teams. ([Netflix TechBlog][3])

After mitigation, the manager ensures learning and follow-through. They do not turn the postmortem into a trial.

---

# 5. How an SRE manager plans

A good SRE manager plans across multiple time horizons.

## Immediate: hours to days

Focus:

* Active incidents
* Dangerous changes
* Capacity emergencies
* Security or data-integrity concerns
* Staff exhaustion
* Immediate mitigations

Questions:

* What can harm users now?
* What is the safest mitigation?
* Is the blast radius contained?
* Who has decision authority?
* What change should be stopped?

The objective is **stability and controlled recovery**.

---

## Short term: one to six weeks

Focus:

* Incident corrective actions
* Alert cleanup
* Capacity increases
* Runbook gaps
* Repeated operational failures
* Upcoming launches
* On-call readiness
* Small automation projects

Example:

```text
Problem:
Seven production pages in two weeks caused by API server saturation.

Short-term work:
- Correct alert thresholds.
- Add request-priority visibility.
- Increase headroom.
- Create overload runbook.
- Add temporary admission limits.
- Test mitigation during a game day.
```

The objective is **risk reduction and operational control**.

---

## Quarterly: approximately three months

Focus:

* Reliability objectives
* Major toil-reduction projects
* Capacity programmes
* Architecture improvements
* Dependency resilience
* Service onboarding
* Staffing and skill gaps
* On-call sustainability
* Reliability debt

A good quarterly plan starts with evidence:

```text
Inputs
 ├── SLO performance
 ├── Error-budget consumption
 ├── Incident themes
 ├── Page and ticket volume
 ├── Capacity forecast
 ├── Launch calendar
 ├── Dependency risk
 ├── Security requirements
 ├── Engineer feedback
 └── Business priorities
```

It then produces a small number of outcomes.

### Weak quarterly objective

> Improve Kubernetes platform reliability.

### Better quarterly objective

> Reduce failed production deployments caused by platform dependencies from 0.6% to below 0.15%, while cutting after-hours pages from 14 to fewer than 5 per week.

### Possible key results

* Introduce dependency-level timeout and retry policies.
* Move three global admission dependencies to regional isolation.
* Add deployment success SLO.
* Automate rollback for two failure categories.
* Remove the five noisiest non-actionable alerts.
* Complete regional failure testing.

Meta has described SLOs as useful in both daily reliability conversations and longer-term strategic planning. ([Engineering at Meta][12])

---

## Medium term: six to twelve months

Focus:

* Fundamental architecture changes
* Multi-region strategy
* Platform self-service
* Fleet automation
* Disaster recovery
* Elimination of operational bottlenecks
* Standardisation
* Service ownership model
* Team topology
* Staff-engineer development

Example:

```text
Current state:
Cluster upgrades require a central team to coordinate every environment.

Six-to-twelve-month target:
- Upgrades are automatically scheduled.
- Risk tiers determine rollout speed.
- Health validation stops bad waves.
- Application teams can request controlled exceptions.
- Rollback and pause are automated.
- Central engineers intervene only for abnormal conditions.
```

The objective is **changing the operating model**, not merely improving individual runbooks.

---

## Long term: one to three years

Focus:

* Where reliability ownership should live
* Platform versus product boundaries
* Architecture that removes entire failure categories
* Global capacity and regional strategy
* Skills the team will require
* Technology retirement
* Control-plane simplification
* Organisational scalability
* Cost-efficiency at future scale

Long-term questions include:

* What breaks when traffic grows 10×?
* What currently requires rare expert knowledge?
* Which systems have unacceptable blast radius?
* Which central components create organisation-wide dependency?
* What should become self-service?
* What should be fully automatic?
* Which service should be retired rather than improved?
* Where is the organisation accumulating irreversible complexity?

The manager’s long-term output may be a strategy such as:

> “Move from a centrally operated Kubernetes fleet to a policy-driven platform where application teams consume paved-road capabilities and the platform team operates shared control planes, standards and exception mechanisms.”

---

# 6. How they allocate team capacity

A weak manager plans 100% of engineer time and then acts surprised when incidents interrupt projects.

A stronger manager explicitly models operational demand.

An illustrative quarterly allocation might be:

```text
50%  Strategic engineering projects
20%  Reliability debt and incident follow-ups
10%  Toil elimination
10%  On-call and expected operational work
10%  Interrupt and uncertainty buffer
```

This is not a universal FAANG formula. The point is that:

* Operational work must be measured.
* Interrupt capacity must exist.
* Reliability work cannot depend on engineers having spare time.
* Planned work must be reduced when operational load rises.
* A team persistently spending most of its time on operations requires structural intervention.

The manager should revisit allocation based on page volume, incident severity, launch demand and team maturity.

---

# 7. How a good SRE manager thinks

## 7.1 Start with user harm, not infrastructure symptoms

Bad question:

> “Why did CPU reach 95%?”

Better questions:

> “Did users experience increased latency or failed requests?”

> “Was the CPU increase the cause, a symptom or harmless utilisation?”

SRE should measure systems from the user’s perspective where possible. Google’s service best practices explicitly recommend defining SLOs like a user. ([Google SRE][13])

---

## 7.2 Reliability is a trade-off, not a maximum

An inexperienced manager may try to make everything “five nines.”

A good manager asks:

* What does additional reliability cost?
* What user value does it provide?
* Does the dependency support the same target?
* Would this investment be better spent elsewhere?
* Is 99.99% required, or would 99.9% be sufficient?
* Are we slowing delivery to protect reliability users do not require?

Pursuing 100% reliability can create excessive cost, complexity and conservatism. Error budgets make acceptable risk explicit. ([Google SRE][14])

---

## 7.3 Treat failures as expected

Bad thinking:

> “How do we ensure this never fails?”

Better thinking:

> “When this fails, how do we detect it, contain it, degrade safely, recover and learn?”

The manager looks for:

* Isolation
* Redundancy
* Backpressure
* Load shedding
* Circuit breaking
* Timeouts
* Idempotency
* Progressive rollout
* Automated rollback
* Graceful degradation
* Regional evacuation
* Restore testing

A system is not resilient because it has never failed. It is resilient because it can absorb and recover from failure.

---

## 7.4 Optimise for blast radius

The probability of failure is only one dimension.

A change may fail infrequently but affect every service, region or customer.

A good manager asks:

```text
Risk ≈ Probability × Impact × Exposure duration
```

For a Kubernetes platform:

* A bad node image may affect one node pool.
* A bad admission-policy rollout may block deployments across hundreds of clusters.
* A faulty global controller may mutate every workload.
* A credential issue may prevent every cluster from reconciling.
* A regional dependency may unexpectedly become a global dependency.

The manager pays disproportionate attention to systems with broad blast radius.

---

## 7.5 Prefer mechanisms over instructions

Bad corrective action:

> “Engineers must remember to verify rollback before deployment.”

Good corrective action:

> “The deployment pipeline blocks release unless rollback validation succeeds.”

Bad:

> “Check capacity before major events.”

Good:

> “Event forecasts automatically generate capacity simulations and block launch if projected headroom is insufficient.”

Rules and reminders decay. Mechanisms make the safe behaviour the default.

---

## 7.6 Treat on-call pain as product feedback

A page is not merely an interruption. It is evidence that the system could not safely manage a condition by itself.

The manager asks:

* Was this page actionable?
* Could the condition self-heal?
* Could it degrade safely?
* Was the alert based on user impact?
* Did the responder have enough context?
* Was expert knowledge required?
* Has this happened before?
* Why was the previous fix insufficient?

Meta has reported using automation to standardise root-cause investigations and reduce repetitive on-call effort, illustrating the broader principle of turning operational expertise into reusable systems. ([Engineering at Meta][15])

---

## 7.7 Optimise the sociotechnical system

Many incidents are not purely code failures.

They may involve:

* Unclear ownership
* Conflicting incentives
* Missing review
* Excessive permissions
* Poor handoff
* Organisational boundaries
* Incomplete training
* Hidden dependencies
* Deadline pressure
* Alert fatigue
* Lack of recovery practice

A good manager investigates the conditions around the technical failure, not only the faulty line of code.

---

## 7.8 Protect cognitive capacity

Engineers working under constant interruption make poorer decisions and produce less durable work.

The manager protects:

* Focus time
* Recovery after incidents
* Reasonable on-call load
* Clear escalation
* Small incident-response groups
* Knowledge sharing
* Predictable priorities
* Time for deep engineering

People health is part of reliability engineering, not a separate human-resources concern.

---

# 8. Good versus average SRE manager

| Area                 | Average manager                 | Good manager                                         |
| -------------------- | ------------------------------- | ---------------------------------------------------- |
| Reliability          | Tracks uptime                   | Connects reliability to critical user journeys       |
| Incidents            | Measures incident count         | Studies impact, detection, mitigation and recurrence |
| On-call              | Ensures rotation coverage       | Ensures the rotation is sustainable and useful       |
| Alerts               | Accepts inherited alerts        | Demands actionable, user-relevant paging             |
| Toil                 | Distributes it fairly           | Engineers it away                                    |
| Planning             | Collects project requests       | Allocates capacity according to risk and leverage    |
| Technical leadership | Makes every decision            | Builds technical leaders and decision mechanisms     |
| Product relationship | Blocks unsafe launches          | Creates conditions for safe launches                 |
| Postmortems          | Finds the mistake               | Finds the system conditions that enabled impact      |
| Performance          | Rewards visible firefighting    | Rewards prevention and durable engineering           |
| Staffing             | Requests headcount for workload | Reduces demand, then justifies residual staffing     |
| Metrics              | Reports activity                | Reports outcomes                                     |
| Documentation        | Requests more runbooks          | Asks why the task is manual                          |
| Ownership            | Lets SRE absorb everything      | Defines and enforces explicit ownership              |
| Long-term thinking   | Improves current operations     | Changes the operating model                          |
| Communication        | Reports technical details       | Explains user impact, risk, options and decisions    |

---

# 9. Detailed good-and-bad examples

## Example 1: Noisy alerts

### Situation

The team receives 40 alerts per week. Only five require action.

### Bad manager

* Tells engineers to respond faster.
* Adds more engineers to the rotation.
* Writes more runbooks.
* Tracks mean time to acknowledge.
* Praises people who handle many alerts.

The system continues producing noise.

### Good manager

First categorises alerts:

```text
40 weekly alerts
├── 18 automatically recover without action
├── 9 are duplicate symptoms
├── 8 can wait until working hours
└── 5 require immediate action
```

Then changes the system:

* Removes non-actionable pages.
* Groups duplicate alerts.
* Converts nonurgent conditions into tickets.
* Alerts on SLO burn rather than every low-level symptom.
* Adds automated mitigation.
* Reviews page usefulness after every shift.

Result:

* Five meaningful pages remain.
* Engineers trust alerts.
* Response quality improves.
* Project time is restored.

---

## Example 2: Database saturation

### Situation

The database reaches connection limits during traffic spikes.

### Bad manager

* Adds an alert at 80%.
* Creates a restart runbook.
* Increases the connection limit.
* Assigns an engineer to watch the dashboard during major events.

This may delay the next incident without resolving the failure mode.

### Good manager

Asks:

1. Is demand legitimate or caused by retries?
2. Which callers consume connections?
3. Is the service applying backpressure?
4. Are connection pools correctly bounded?
5. Can low-priority traffic be shed?
6. What happens when the database refuses requests?
7. Is failover tested?
8. Is the capacity forecast based on peak concurrency or averages?

Possible engineering response:

* Bound client pools.
* Make retries exponential and budgeted.
* Add per-client quotas.
* Shed noncritical requests.
* Add admission control.
* Cache selected reads where safe.
* Improve query performance.
* Test database failover under load.
* Alert on user-visible saturation and budget burn.

The difference is that the good manager treats saturation as a system-design problem, not a dashboard problem.

---

## Example 3: Major product launch

### Bad manager

SRE is invited one week before launch.

The manager asks:

* Are dashboards ready?
* Is there a runbook?
* Who is on-call?

The launch proceeds after informal approval.

### Good manager

SRE engages during design.

Six weeks before launch:

* Define critical journeys.
* Establish SLOs.
* Estimate load.
* Identify dependencies.
* Model failure modes.
* Define degradation.
* Review data-loss risk.

Three weeks before launch:

* Run load tests.
* Validate alerting.
* Test rollback.
* Confirm capacity.
* Conduct dependency failure tests.

Launch week:

* Canary in one environment.
* Expand gradually.
* Monitor user-facing indicators.
* Enforce stop conditions.
* Keep rollback authority clear.

After launch:

* Review actual load and failures.
* Remove temporary procedures.
* Update the long-term reliability roadmap.

The good manager does not become the final approval gate. They create a repeatable launch-safety mechanism.

---

## Example 4: Error budget is exhausted

### Bad response A: Ignore it

> “The business launch is important. Keep shipping.”

Then the SLO has no decision-making value.

### Bad response B: Freeze everything indiscriminately

> “No production changes until next quarter.”

This may block changes that reduce risk.

### Good response

Apply a previously agreed policy:

* Pause high-risk feature rollouts.
* Continue low-risk fixes and reliability improvements.
* Identify the dominant budget-consuming events.
* Add stronger canary requirements.
* Increase leadership review for exceptions.
* Resume normal release velocity once burn stabilises.

The key is that the response was agreed before the conflict.

---

## Example 5: Repeated certificate expiry incident

### Bad manager

Postmortem action:

* Update the renewal runbook.
* Add two calendar reminders.
* Ask the on-call engineer to check monthly.

### Good manager

Finds the systemic gaps:

* Inventory was incomplete.
* Ownership was unclear.
* Expiry monitoring covered only public certificates.
* Renewal succeeded but deployment failed.
* No end-to-end validation existed.

Corrective actions:

* Central certificate inventory
* Automated renewal
* Deployment verification
* Expiry SLO
* Ownership metadata
* Progressive rollout
* Alert only when automation cannot recover
* Removal of manually managed certificates

---

## Example 6: Team burnout

### Bad manager

* Thanks the team for resilience.
* Organises a team lunch.
* Adds another person to the rotation.
* Tells engineers to take time off.
* Keeps the same service scope and alert load.

### Good manager

Quantifies the problem:

```text
Last 30 days:
- 72 after-hours pages
- 31 non-actionable
- 19 repeated failure mode
- 4 engineers carrying 80% of escalations
- 2 planned projects repeatedly interrupted
```

Then:

* Removes noisy pages.
* Assigns a project to eliminate the repeated failure.
* Trains additional responders.
* Transfers inappropriate ownership.
* Reduces quarterly commitments.
* Gives responders recovery time.
* Escalates unsafe service conditions.
* Reviews whether the team should continue supporting the service.

The team lunch may still happen, but it is not presented as the solution.

---

## Example 7: Strong engineer becomes indispensable

### Bad manager

The manager keeps assigning the most difficult incidents to that engineer because they are fastest.

Result:

* Other engineers do not develop.
* The expert burns out.
* The team has a single point of failure.
* The manager mistakes dependency for high performance.

### Good manager

* Makes the expert a mentor rather than default responder.
* Assigns paired investigations.
* Rotates subsystem ownership.
* Creates reverse-engineering exercises.
* Documents mental models, not merely procedures.
* Measures escalation concentration.
* Builds at least two owners for every critical area.

Google has described training practices involving design overviews, production deep dives and hands-on exercises before teams assume operational ownership. ([Google SRE][8])

---

# 10. Example: SRE manager for a large Kubernetes platform

Assume the team manages:

* Hundreds of Kubernetes clusters
* Thousands of services
* Cluster lifecycle
* Upgrades
* Networking
* Add-ons
* Deployment systems
* Policy
* Observability
* Incident response

## Weak operating model

The manager measures:

* Number of clusters upgraded
* Tickets closed
* Incidents handled
* Automation scripts delivered
* Number of dashboards

The team spends most of its time:

* Manually fixing failed upgrades
* Approving exceptions
* Investigating application problems
* Restarting components
* Coordinating changes
* Answering repeated support questions

This looks busy but does not scale.

## Strong operating model

The manager defines critical user journeys:

1. Create a cluster.
2. Upgrade a cluster safely.
3. Deploy a workload.
4. Scale a workload.
5. Recover from node or zonal failure.
6. Apply policy consistently.
7. Access logs and metrics.
8. Roll back a dangerous platform change.

They then establish indicators:

```text
Cluster provisioning success
Cluster provisioning duration
Upgrade success without manual intervention
Workload deployment success
Control-plane availability
Admission latency
DNS success
Network connectivity
Node replacement time
Policy propagation delay
Fleet configuration convergence
Recovery time
Manual intervention rate
```

## Example quarterly plan

### Objective 1: Reduce manual fleet operations

Key results:

* Increase automatic cluster-upgrade completion from 82% to 96%.
* Reduce upgrade-related engineer interventions by 60%.
* Automatically classify the top five upgrade failure modes.
* Add automatic pause and rollback at the rollout-wave level.

### Objective 2: Reduce global blast radius

Key results:

* Remove two synchronous global dependencies from cluster reconciliation.
* Introduce regional isolation for admission services.
* Add rate limiting to global controllers.
* Test loss of the central management plane.

### Objective 3: Improve developer experience without reducing safety

Key results:

* Reduce median cluster-provisioning time from six hours to two.
* Provide self-service exception requests with automated policy checks.
* Show actionable failure reasons in the platform UI.
* Reduce platform support tickets by 30%.

### Objective 4: Improve on-call sustainability

Key results:

* Reduce after-hours pages from 15 to fewer than six per week.
* Eliminate three recurring page classes.
* Train every engineer on at least two major subsystems.
* Ensure no subsystem depends on one responder.

This plan connects reliability, scalability, developer productivity and team health.

---

# 11. What separates a good manager from an exceptional one

A good SRE manager can operate their own team well.

An exceptional SRE manager improves the reliability behaviour of the broader engineering organisation.

They create reusable mechanisms such as:

* Company-wide SLO standards
* Production-readiness frameworks
* Safe rollout platforms
* Incident-management tooling
* Dependency-risk reviews
* Automated capacity planning
* Standard service templates
* Common game-day programmes
* Reliability scorecards
* Self-service operational platforms
* Better incentives for service ownership

They move from:

> “My team handled this problem.”

to:

> “The platform and engineering model now prevent hundreds of teams from encountering this problem.”

This is leverage.

---

# 12. A practical decision framework

A strong SRE manager can evaluate most problems through the following sequence.

## Step 1: Define the user impact

* Who is affected?
* Which journey fails?
* What is the severity?
* How long does it last?

## Step 2: Quantify the risk

* How often can it occur?
* What is the blast radius?
* How quickly is it detected?
* How quickly can it be mitigated?
* Could data be lost or corrupted?

## Step 3: Identify the dominant failure mechanisms

* Change failure?
* Capacity exhaustion?
* Dependency failure?
* Configuration drift?
* Human error?
* Missing isolation?
* Unbounded retries?
* Unknown ownership?

## Step 4: Choose the appropriate control

In descending order of preference:

1. Remove the failure mode.
2. Prevent it.
3. Contain it.
4. Detect it.
5. Automatically mitigate it.
6. Help a human mitigate it.
7. Document manual recovery.

Many average organisations jump directly to step 7.

## Step 5: Assign ownership and deadlines

Every important action needs:

* One accountable owner
* Expected risk reduction
* Due date
* Verification method
* Status review

## Step 6: Verify that the risk actually declined

Completing a ticket is not proof.

Measure:

* Did pages decline?
* Did SLO performance improve?
* Did recovery become faster?
* Did manual intervention decline?
* Was the failure tested?
* Is the blast radius smaller?

---

# 13. The simplest mental model

A good SRE manager repeatedly moves the organisation through this progression:

```text
Reactive
  Humans repeatedly repair production
        ↓
Documented
  Humans follow consistent procedures
        ↓
Automated
  Software performs known procedures
        ↓
Self-healing
  Systems detect and mitigate known failures
        ↓
Resilient
  Architecture contains failures and degrades safely
        ↓
Adaptive
  Reliability mechanisms evolve with load, products and risks
```

An average manager makes the reactive stage more organised.

A good manager moves the team upward.

---

# 14. Questions a good SRE manager asks regularly

### About production

* What is currently most likely to cause significant user harm?
* Which system has the largest blast radius?
* Which dependency do we trust without evidence?
* Which recovery procedure has never been tested?
* Which successful metric is hiding a bad user experience?

### About operations

* Why did this require a human?
* Which page should never happen again?
* What percentage of our work is interrupt-driven?
* Where is operational knowledge concentrated?
* Which team owns the underlying cause?

### About planning

* What are we deliberately not doing?
* How much unplanned work can the team absorb?
* Which project removes the most future work?
* Are we optimising a system that should be retired?
* Is the roadmap based on risk or stakeholder volume?

### About people

* Who is overloaded?
* Who is not receiving meaningful ownership?
* Who is becoming a single point of failure?
* Are engineers recognised for prevention?
* Does the team have uninterrupted engineering time?

### About leadership

* Are reliability decisions based on data or hierarchy?
* Are product teams and SRE using the same definition of success?
* Have we made safe behaviour easy?
* Are exceptions visible and temporary?
* Are we relying on heroics?

---

## Final distinction

An average SRE manager keeps the current operating model functioning.

A good SRE manager continuously redesigns the operating model so that:

* fewer incidents occur,
* failures have smaller impact,
* recovery is faster,
* changes are safer,
* engineers receive fewer interruptions,
* product teams retain ownership,
* operational knowledge is distributed,
* and reliability improves without requiring the team to grow linearly.

Their deepest question is not:

> “How do we operate this system?”

It is:

> **“How should we engineer the system, the platform and the organisation so that operating it becomes progressively easier and safer?”**

[1]: https://sre.google/sre-book/introduction/?utm_source=chatgpt.com "Google SRE - IT Service Management: Automate Operations"
[2]: https://engineering.fb.com/category/production-engineering/?utm_source=chatgpt.com "Production Engineering Archives - Engineering at Meta"
[3]: https://netflixtechblog.com/keeping-customers-streaming-the-centralized-site-reliability-practice-at-netflix-205cc37aa9fb?utm_source=chatgpt.com "The Centralized Site Reliability Practice at Netflix"
[4]: https://sre.google/resources/practices-and-processes/product-focused-reliability-for-sre/?utm_source=chatgpt.com "Product SRE, improving reliability of services"
[5]: https://sre.google/workbook/error-budget-policy/?utm_source=chatgpt.com "Error Budget Policy for Service Reliability"
[6]: https://sre.google/sre-book/dealing-with-interrupts/?utm_source=chatgpt.com "Chapter 29 - Dealing with Interrupts"
[7]: https://sre.google/workbook/eliminating-toil/?utm_source=chatgpt.com "Operational Efficiency: Eliminating Toil"
[8]: https://sre.google/sre-book/evolving-sre-engagement-model/?utm_source=chatgpt.com "Production Readiness Review: Engagement Insight"
[9]: https://aws.amazon.com/builders-library/building-dashboards-for-operational-visibility/?utm_source=chatgpt.com "Building dashboards for operational visibility"
[10]: https://docs.aws.amazon.com/wellarchitected/latest/operational-excellence-pillar/ops_mit_deploy_risks_deploy_mgmt_sys.html?utm_source=chatgpt.com "OPS06-BP03 Employ safe deployment strategies"
[11]: https://sre.google/sre-book/managing-incidents/?utm_source=chatgpt.com "Incident Management: Key to Restore Operations"
[12]: https://engineering.fb.com/2021/12/13/production-engineering/slick/?utm_source=chatgpt.com "SLICK: Adopting SLOs for improved reliability"
[13]: https://sre.google/sre-book/service-best-practices/?utm_source=chatgpt.com "Google SRE: Production Services Best Practices"
[14]: https://sre.google/sre-book/service-level-objectives/?utm_source=chatgpt.com "Defining slo: service level objective meaning"
[15]: https://engineering.fb.com/2025/12/19/data-infrastructure/drp-metas-root-cause-analysis-platform-at-scale/?utm_source=chatgpt.com "DrP: Meta's Root Cause Analysis Platform at Scale"
