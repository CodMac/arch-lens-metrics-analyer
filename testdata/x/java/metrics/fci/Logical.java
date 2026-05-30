package metrics.fci;

public class Logical {
    // CC = 1 (base) + 1 (if) + 1 (&&) + 1 (||) = 4
    public boolean testLogical(boolean a, boolean b, boolean c) {
        if (a && b || c) {
            return true;
        }
        return false;
    }

    // CC = 1 (base) + 1 (if) + 2 (&&) = 4
    public void nestedLogical(boolean a, boolean b, boolean c) {
        if (a && (b && c)) {
            // ...
        }
    }
}
// Expected FCI = 8 (4 + 4)
