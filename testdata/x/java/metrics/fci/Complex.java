package metrics.fci;

public class Complex {
    // CC = 1 + 1 (if) + 1 (&&) = 3
    public void m1(int a, int b) {
        if (a > 0 && b > 0) {
            System.out.println("Both positive");
        }
    }

    // CC = 1 + 3 (case) = 4
    public void m2(int type) {
        switch (type) {
            case 1:
            case 2:
            case 3:
                System.out.println("Type 1, 2 or 3");
                break;
            default:
                System.out.println("Other");
        }
    }

    // CC = 1 + 1 (for) + 1 (if) = 3
    public void m3(int[] arr) {
        for (int i : arr) {
            if (i % 2 == 0) {
                System.out.println(i);
            }
        }
    }
}
// Expected FCI = 3 + 4 + 3 = 10
