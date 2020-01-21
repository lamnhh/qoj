#include <bits/stdc++.h>
using namespace std;

typedef long long LL;
typedef pair<int, int> II;

const int MAXN = (int)2e5 + 10;
int n, a[MAXN], b[MAXN], c[MAXN];
char str[MAXN];

void read(int a[]) {
    scanf("%s", str + 1);
    for (int i = 1; i <= n; ++i) {
        a[n - i] = str[i] - 'a';
    }
}

int main() {
    scanf("%d", &n);
    read(a);
    read(b);

    n += 1;
    int carry = 0;
    for (int i = 0; i < n; ++i) {
        c[i] = (a[i] + b[i] + carry) % 26;
        carry = (a[i] + b[i] + carry) / 26;
    }
    assert(carry == 0);

    carry = 0;
    for (int i = n - 1; i >= 0; --i) {
        carry = carry * 26 + c[i];
        if (i < n - 1) {
            putchar('a' + (carry / 2));
        }
        carry %= 2;
    }
    puts("");
    return 0;
}