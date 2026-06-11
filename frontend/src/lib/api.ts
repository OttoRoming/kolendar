export class APIError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'APIError';
	}
}

export class AuthError extends Error {
	constructor() {
		super('Unauthorized');
		this.name = 'AuthError';
	}
}

async function post<T>(url: string, body: unknown): Promise<T> {
	let response: Response;
	try {
		response = await fetch(url, {
			method: 'POST',
			credentials: 'same-origin',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(body)
		});
	} catch {
		throw new APIError('Network error');
	}

	if (response.status === 401) {
		throw new AuthError();
	}
	if (response.status === 502) {
		throw new APIError('Server is down');
	}
	if (!response.ok) {
		const data = await response.json();
		const message = data.message || 'An error occurred';
		throw new APIError(message);
	}

	return response.json();
}

interface UserRequest {
	username: string;
	password: string;
}

interface UserResponse {
	id: string;
	username: string;
	token: string;
}

async function signup(body: UserRequest): Promise<UserResponse> {
	return post('/api/users', body);
}

async function login(body: UserRequest): Promise<UserResponse> {
	return post('/api/users/login', body);
}

const api = {
	signup,
	login
};

export default api;
