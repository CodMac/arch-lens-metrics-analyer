package metrics.fci;

import java.util.function.Consumer;

public class LambdaAndAnonymous {
    // CC = 1 (base) + 1 (if in lambda) + 1 (if in anonymous) = 3
    // According to md: lambda's CC is added to host method.
    public void test() {
        Consumer<String> c = s -> {
            if (s != null) { // +1
                System.out.println(s);
            }
        };

        new Runnable() {
            @Override
            public void run() {
                if (true) { // +1
                    // ...
                }
            }
        }.run();
    }
}
// Expected FCI = 3
