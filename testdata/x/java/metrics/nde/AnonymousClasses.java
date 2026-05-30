package metrics.nde;

public class AnonymousClasses {
    public void test() {
        Runnable r = new Runnable() {
            @Override
            public void run() {}
        };
    }
}
// Expected NDE = 1 (Spec)
