package metrics.fci;

public class DeepNesting {
    // CC = 1 (base) + 1 (if) + 1 (while) + 1 (for) + 1 (if) = 5
    public void deep() {
        if (true) {
            while (true) {
                for (int i = 0; i < 10; i++) {
                    if (i == 5) {
                        break;
                    }
                }
                break;
            }
        }
    }
}
// Expected FCI = 5
