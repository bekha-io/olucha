# Olucha

<img src="https://i.ibb.co/Tqd7586H/logo.png" alt="logo">

> Warning! The product still is in active development phase and doesn't provide the mentioned functionality yet. For the detailed information, please message iskandarzoda@outlook.com

**Olucha** - is an open source workflow orchestration engine that provides an advanced JSON DSL to define stateful workflows and steps. It also provides a simplified UI for building the workflows for those who doesn't want to specify it via DSL

# Definitions

## Workflow

Workflow is a main entity that represents a sequence of steps that need to be executed to achieve a specific business goal. Each workflow has a unique identifier and a version. Workflow can contain multiple steps and forms that are used to collect data from users. Steps can be of different types, such as human tasks, automated tasks, or decision points. Forms are used to define the structure of data that needs to be collected from users during workflow execution.

## Step

Step is a single action within a workflow. Each step has a unique identifier and can be one of the following types:

- **humanTask** - a step that requires human interaction (e.g. filling out a form)
- **systemTask** - an automated step performed by the system
- **condition** - a decision step based on certain conditions

Each step can contain:

- ID - step's unique identifier
- Label - human readable name for the step
- Description - detailed description of the action

### System Task

System task is an automated step that is executed by the system without human interaction. It can be used to perform various operations such as:

- Making external API calls (HTTP(S) only as of now)

Each system task has the following properties:

- Type - type of operation to perform (e.g. "http", "database", "notification")
- Timeout - optional maximum execution time, after which the workflow considered failed if no response returned by that time

### Human Task

Human task is a step that requires human interaction. It can be used to collect data from users, approve or reject requests, or perform any other action that requires human input. Human task must have a form associated with it that defines the structure of data that needs to be collected from users.

Each human task has the following properties:

- Form ID - identifier of the form that needs to be filled out
- RBAC - specific roles or permissions that are allowed to perform this specific task. Roles ( permissions) are passed to the executor ad hoc

### Condition

Condition is a decision step that allows branching workflow execution based on certain criteria. It evaluates a boolean expression and routes the workflow to different paths depending on the result.

Each condition has:

- Expression - boolean expression to evaluate (e.g. "amount > 1000")
- True path - step to execute if condition is true
- False path - step to execute if condition is false
- Default path - optional step to execute if expression evaluation fails
