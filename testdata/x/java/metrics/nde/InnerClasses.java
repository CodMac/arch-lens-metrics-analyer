package metrics.nde;

public class InnerClasses {
    public class Inner1 {}
    private static class Inner2 {
        class DeepInner {}
    }
}
// Expected NDE = 1 (Spec) or 4 (Current Impl counting all class-likes in file)
