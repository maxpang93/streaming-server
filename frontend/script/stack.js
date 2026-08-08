export default class Stack {
    constructor(stackName) {
        this.stackName = stackName;
    }

    getStack() {
        return JSON.parse(sessionStorage.getItem(this.stackName)) || [];
    }

    updateStack(stack) {
        sessionStorage.setItem(this.stackName, JSON.stringify(stack));
    }

    pop() {
        const stack = this.getStack();
        if (stack.length < 1) {
            return null;
        }
        const lastItem = stack.pop();
        this.updateStack(stack)
        return lastItem;
    }

    push(item) {
        const stack = this.getStack();
        stack.push(item);
        this.updateStack(stack)
    }

    reset() {
        this.updateStack([])
    }
}
