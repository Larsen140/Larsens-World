/* Larsens World Starter: öffnet die App von GitHub Pages in Edge (sonst Chrome) im App-Modus.
   Die App selbst steckt nicht in der exe, deshalb kommen Updates automatisch. */
#include <windows.h>
#include <shellapi.h>
#define APP_URL L"https://larsen140.github.io/Larsens-World/?app=1"

int WINAPI wWinMain(HINSTANCE a, HINSTANCE b, LPWSTR c, int d) {
    const wchar_t *args = L"--app=" APP_URL L" --window-size=1440,900";
    HINSTANCE r = ShellExecuteW(NULL, L"open", L"msedge.exe", args, NULL, SW_SHOWNORMAL);
    if ((INT_PTR)r <= 32) r = ShellExecuteW(NULL, L"open", L"chrome.exe", args, NULL, SW_SHOWNORMAL);
    if ((INT_PTR)r <= 32) r = ShellExecuteW(NULL, L"open", APP_URL, NULL, NULL, SW_SHOWNORMAL);
    return (INT_PTR)r <= 32;
}
