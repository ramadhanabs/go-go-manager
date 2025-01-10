export function randomEmail() {
    const domains = ['example.com', 'test.com', 'demo.com'];
    const username = Math.random().toString(36).substring(2, 8);
    const domain = domains[Math.floor(Math.random() * domains.length)];
    return `${username}@${domain}`;
}

export function randomPassword() {
    return Math.random().toString(36).substring(2, 10);
}

export function randomAction() {
    const actions = ['login', 'signup'];
    return actions[Math.floor(Math.random() * actions.length)];
}

export function logResult(description, success) {
    const green = '\x1b[32m';
    const red = '\x1b[31m';
    const reset = '\x1b[0m';
    console.log(`${success ? green : red}${description}${reset}`);
}