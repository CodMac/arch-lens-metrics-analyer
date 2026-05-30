package metrics.fci;

public class InnerClasses {
    // CC = 1
    public void outer() {}

    public class Inner {
        // CC = 1 + 1 (if) = 2
        public void innerMethod() {
            if (true) {}
        }

        public class DeepInner {
            // CC = 1 + 1 (while) = 2
            public void deepMethod() {
                while(true) break;
            }
        }
    }
    
    // Static inner class
    public static class StaticInner {
        // CC = 1
        public void staticInnerMethod() {}
    }
}
// Expected FCI = 1 + 2 + 2 + 1 = 6
