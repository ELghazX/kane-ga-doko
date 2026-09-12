import { browser } from '$app/environment';

export const authState = $state({
	token: browser ? localStorage.getItem('token') : null,
	userId: browser ? Number(localStorage.getItem('user_id')) : null,
	username: browser ? localStorage.getItem('username') : null,
	isAdmin: browser ? localStorage.getItem('is_admin') === 'true' : false,
});

export function login(data: { token: string, user_id: number, username: string, is_admin: boolean }) {
	authState.token = data.token;
	authState.userId = data.user_id;
	authState.username = data.username;
	authState.isAdmin = data.is_admin;
	
	if (browser) {
		localStorage.setItem('token', data.token);
		localStorage.setItem('user_id', data.user_id.toString());
		localStorage.setItem('username', data.username);
		localStorage.setItem('is_admin', data.is_admin.toString());
	}
}

export function logout() {
	authState.token = null;
	authState.userId = null;
	authState.username = null;
	authState.isAdmin = false;
	
	if (browser) {
		localStorage.removeItem('token');
		localStorage.removeItem('user_id');
		localStorage.removeItem('username');
		localStorage.removeItem('is_admin');
	}
}

export const API_URL = "/api";
