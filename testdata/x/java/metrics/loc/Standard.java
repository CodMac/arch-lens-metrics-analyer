package metrics.loc;

import java.util.List;

/**
 * Standard class with comments and blank lines.
 */
public class Standard {
    // A comment
    private int x = 1;

    public void test() {
        /*
         * Multi-line
         * comment
         */
        if (true) {
            System.out.println("Hello");
        }
    }
}

// L1: package
// L2: (empty)
// L3: import
// L4: (empty)
// L5: /**
// L6:  * Standard...
// L7:  */
// L8: public class Standard {
// L9:     // A comment
// L10:    private int x = 1;
// L11:
// L12:    public void test() {
// L13:        /*
// L14:         * Multi-line
// L15:         * comment
// L16:         */
// L17:        if (true) {
// L18:            System.out.println("Hello");
// L19:        }
// L20:    }
// L21: }
// Total: 21 lines.
// Empty: L2, L4, L11 (3)
// Comment: L5, L6, L7, L9, L13, L14, L15, L16 (8)
// Logical: 21 - 3 - 8 = 10
// Lines: 1, 3, 8, 10, 12, 17, 18, 19, 20, 21.

