package metrics.fci;

public class ControlFlow {
    // CC = 1 (base) + 1 (if) + 1 (else if) + 1 (for) + 1 (while) + 1 (do-while) + 2 (case) + 1 (catch) = 9
    public void complex(int x) {
        if (x > 0) {
            // ...
        } else if (x < 0) { // +1
            // ...
        }

        for (int i = 0; i < 10; i++) { // +1
            // ...
        }

        while (x < 100) { // +1
            x++;
        }

        do { // +1
            x--;
        } while (x > 0);

        switch (x) {
            case 1: // +1
                break;
            case 2: // +1
                break;
            default:
                break;
        }

        try {
            x = 1 / 0;
        } catch (ArithmeticException e) { // +1
            x = 0;
        } finally {
            x = -1;
        }
    }
}
// Expected FCI = 9
