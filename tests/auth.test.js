import http from 'k6/http';
import { check } from 'k6';
import { logResult, randomEmail, randomPassword } from "./utils/utils.js"

export default function () {
    const url = 'http://localhost:8888/api/v1/auth';
    let payload = {}
    let params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };
    let res = null

    // to store name and password after registration
    let email = ""
    let password = ""

    console.log("--------- START OF REGISTER SCENARIO ---------");
    console.log("Testing Positive Scenario - REGISTER");
    console.log("> Register with valid email and password");
    email = randomEmail()
    password = randomPassword()

    payload = JSON.stringify({
        email,
        password,
        action: 'signup',
    });

    res = http.post(url, payload, params);

    logResult(
        "Expected 200 status for valid register",
        check(res, { 'is status 200': (r) => r.status === 200 })
    );

    console.log("Testing Negative Scenario - REGISTER");
    console.log("> Register with empty email");
    payload = JSON.stringify({
        email: '',
        password: randomPassword(),
        action: 'signup',
    });
    res = http.post(url, payload, params);

    logResult(
        "Expected 400 status for invalid register",
        check(res, { 'is status 400': (r) => r.status === 400 })
    );

    console.log("> Register with empty password");
    payload = JSON.stringify({
        email: randomEmail(),
        password: '',
        action: 'signup',
    });
    res = http.post(url, payload, params);

    logResult(
        "Expected 400 status for invalid register",
        check(res, { 'is status 400': (r) => r.status === 400 })
    );

    console.log("> Register with invalid email format");
    payload = JSON.stringify({
        email: 'invalid_email',
        password: randomPassword(),
        action: 'signup',
    });
    res = http.post(url, payload, params);

    logResult(
        "Expected 400 status for invalid register",
        check(res, { 'is status 400': (r) => r.status === 400 })
    );

    console.log("> Register with empty password");
    payload = JSON.stringify({
        email: randomEmail(),
        password: '',
        action: 'signup',
    });
    res = http.post(url, payload, params);

    logResult(
        "Expected 400 status for invalid register",
        check(res, { 'is status 400': (r) => r.status === 400 })
    );

    console.log("> Register with short password (char < 8)");
    payload = JSON.stringify({
        email: randomEmail(),
        password: '123456',
        action: 'signup',
    });
    res = http.post(url, payload, params);

    logResult(
        "Expected 400 status for invalid register",
        check(res, { 'is status 400': (r) => r.status === 400 })
    );

    console.log("> Register with long password (char > 32)");
    payload = JSON.stringify({
        email: randomEmail(),
        password: '12312312312312312312312312312312312312312312312321',
        action: 'signup',
    });
    res = http.post(url, payload, params);

    logResult(
        "Expected 400 status for invalid register",
        check(res, { 'is status 400': (r) => r.status === 400 })
    );

    console.log("> Register with long password (char > 32)");
    payload = JSON.stringify({
        email: randomEmail(),
        password: '12312312312312312312312312312312312312312312312321',
        action: 'signup',
    });
    res = http.post(url, payload, params);

    logResult(
        "Expected 400 status for invalid register",
        check(res, { 'is status 400': (r) => r.status === 400 })
    );

    console.log("> Register with existed email");
    payload = JSON.stringify({
        email,
        password: randomPassword(),
        action: 'signup',
    });
    res = http.post(url, payload, params);

    logResult(
        "Expected 409 status for invalid register",
        check(res, { 'is status 409': (r) => r.status === 409 })
    );

    console.log("--------- END OF REGISTER SCENARIO ---------");
    console.log("\n")
    console.log("\n")

    console.log("--------- START OF LOGIN SCENARIO ---------");
    // Positive Scenario - LOGIN
    console.log("Testing Positive Scenario - LOGIN");
    console.log("> Login with valid email and password");
    payload = JSON.stringify({
        email,
        password,
        action: 'login',
    });

    res = http.post(url, payload, params);

    logResult(
        "Expected 200 status for valid login",
        check(res, { 'is status 200': (r) => r.status === 200 })
    );

    // Negative Scenario - LOGIN (Invalid Email)
    console.log("Testing Negative Scenario - LOGIN");
    console.log("> Login with invalid email");
    payload = JSON.stringify({
        email: 'invalid_email',
        password: randomPassword(),
        action: 'login',
    });

    const invalidEmailRes = http.post(url, payload, params);

    logResult(
        "Expected 400 status for invalid email",
        check(invalidEmailRes, { 'is status 400': (r) => r.status === 400 })
    );

    // Negative Scenario - LOGIN (Invalid Password)
    console.log("> Login with invalid password");
    payload = JSON.stringify({
        email: randomEmail(),
        password: '',
        action: 'login',
    });

    const invalidPasswordRes = http.post(url, payload, params);

    logResult(
        "Expected 400 status for invalid password",
        check(invalidPasswordRes, { 'is status 400': (r) => r.status === 400 })
    );

    console.log("--------- END OF LOGIN SCENARIO ---------");
    console.log("\n")
    console.log("\n")
}