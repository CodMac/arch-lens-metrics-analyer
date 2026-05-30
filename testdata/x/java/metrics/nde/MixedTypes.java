package metrics.nde;

public interface MixedTypes {
    void doSomething();
}

enum MyEnum {
    A, B
}

@interface MyAnnotation {
}
// Expected NDE = 3
