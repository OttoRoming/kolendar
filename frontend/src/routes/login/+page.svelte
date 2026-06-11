<script lang="ts">
	import Button from '$lib/Button.svelte';
	import H1 from '$lib/H1.svelte';
	import H2 from '$lib/H2.svelte';
	import Submit from '$lib/Submit.svelte';
	import TextInput from '$lib/TextInput.svelte';
	import { toast } from 'svelte-sonner';
	import api from '$lib/api';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	let signupValue = $state({
		username: '',
		password: '',
		verifyPassword: ''
	});

	let loginValue = $state({
		username: '',
		password: ''
	});

	let formClasses = 'flex flex-col items-center gap-2';

	function signup(e: Event) {
		e.preventDefault();

		if (signupValue.password !== signupValue.verifyPassword) {
			toast.error('Passwords do not match');
			return;
		}

		toast.promise(api.signup(signupValue), {
			loading: 'Signing up...',
			success: () => {
				goto(resolve('/libraries'));
				return 'Signed up successfully!';
			},
			error: (error) => {
				if (error instanceof Error) {
					return error.message;
				}
				return 'An unknown error occurred';
			}
		});
	}

	function login(e: Event) {
		e.preventDefault();

		toast.promise(api.login(loginValue), {
			loading: 'Logging in...',
			success: () => {
				goto(resolve('/libraries'));
				return 'Logged in successfully!';
			},
			error: (error) => {
				if (error instanceof Error) {
					return error.message;
				}
				return 'An unknown error occurred';
			}
		});
	}
</script>

<H1>Authenticate</H1>
<div class="grid grid-cols-2 gap-16">
	<form class={formClasses} onsubmit={signup}>
		<H2>Sign Up</H2>
		<TextInput type="text" label="Username" bind:value={signupValue.username} />
		<TextInput type="password" label="Password" bind:value={signupValue.password} />
		<TextInput type="password" label="Verify Password" bind:value={signupValue.verifyPassword} />
		<Submit value="Register" />
	</form>
	<form class={formClasses} onsubmit={login}>
		<H2>Login</H2>
		<TextInput type="text" label="Username" bind:value={loginValue.username} />
		<TextInput type="password" label="Password" bind:value={loginValue.password} />
		<Submit value="Login" />
	</form>
</div>
