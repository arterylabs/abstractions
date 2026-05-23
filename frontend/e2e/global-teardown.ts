async function globalTeardown() {
  const apiURL =
    process.env.PLAYWRIGHT_API_URL ||
    process.env.NEXT_PUBLIC_API_URL ||
    "http://localhost:1234";

  try {
    const res = await fetch(`${apiURL}/test/cleanup`, { method: "POST" });
    const body = await res.json();
    console.log(`[teardown] cleaned up ${body.deletedCount ?? 0} test user(s)`);
  } catch (err) {
    console.warn("[teardown] cleanup request failed:", err);
  }
}

export default globalTeardown;
