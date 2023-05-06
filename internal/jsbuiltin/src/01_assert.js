{
  function assert(ok, message) {
    if (!ok) {
      throw new AssertionError({ actual: ok, expected: true, message });
    }
  }

  class AssertionError extends Error {
    constructor(options) {
      super();
      this.name = "AssertionError";
      this.actual = options.actual;
      this.expected = options.expected;
      this.message =
        options.message || `expected ${this.actual} == ${this.expected}`;
    }

    toString() {
      return `${this.name}: ${this.message}`;
    }
  }

  jssh.assert = assert;
  jssh.AssertionError = AssertionError;
}
