# Nexus Multi-SDK Implementation Examples

TypeScript (experimental), Python (public preview), and Java (GA) examples for Nexus handler and caller implementations.

## TypeScript (Experimental)

### Service Contract Definition

```typescript
import * as nexus from 'nexus-rpc';
export const orderService = nexus.service('order-service', {
    echo: nexus.operation<EchoInput, EchoOutput>(),
    processOrder: nexus.operation<OrderInput, OrderOutput>(),
});
```

### Handler Implementations

```typescript
import * as temporalNexus from '@temporalio/nexus';
const handlers = {
    echo: async (ctx, input: EchoInput): Promise<EchoOutput> => {
        return { message: input.message };
    },
    processOrder: new temporalNexus.WorkflowRunOperationHandler<OrderInput, OrderOutput>(
        async (ctx, input: OrderInput) => {
            return await temporalNexus.startWorkflow(ctx, processOrderWorkflow, {
                args: [input],
                workflowId: ctx.requestId ?? randomUUID(),
            });
        }
    ),
};
```

### Caller

```typescript
const nexusClient = wf.createNexusClient({ service: orderService, endpoint: 'my-endpoint' });
const result = await nexusClient.executeOperation('processOrder', { id: input.orderId }, { scheduleToCloseTimeout: '10m' });
```

## Python (Public Preview)

### Handler

```python
@nexusrpc.service
class OrderService:
    process_order: nexusrpc.Operation[OrderInput, OrderOutput]
    echo: nexusrpc.Operation[EchoInput, EchoOutput]

@nexusrpc.handler.service_handler(service=OrderService)
class OrderServiceHandler:
    @nexusrpc.handler.sync_operation
    async def echo(self, ctx, input: EchoInput) -> EchoOutput:
        return EchoOutput(message=input.message)

    @temporal_nexus.workflow_run_operation
    async def process_order(self, ctx, input: OrderInput):
        return await temporal_nexus.start_workflow(ctx, ProcessOrderWorkflow, args=[input])
```

## Java (GA)

### Handler

```java
@ServiceImpl(service = OrderNexusService.class)
public class OrderNexusServiceImpl {
    @OperationImpl
    public OperationHandler<EchoInput, EchoOutput> echo() {
        return OperationHandler.sync((ctx, input) -> new EchoOutput(input.getMessage()));
    }
    @OperationImpl
    public OperationHandler<OrderInput, OrderOutput> processOrder() {
        return WorkflowRunOperationHandler.fromWorkflowMethod(
            (ctx, input) -> Workflow.newWorkflowStub(ProcessOrderWorkflow.class)::processOrder);
    }
}
```
