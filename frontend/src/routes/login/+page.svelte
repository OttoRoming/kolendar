<script lang="ts">
	import Button from '$lib/Button.svelte';
	import H1 from '$lib/H1.svelte';
	import H2 from '$lib/H2.svelte';
	import Submit from '$lib/Submit.svelte';
	import TextInput from '$lib/TextInput.svelte';
	import { toast } from 'svelte-sonner';
	import ky, { isHTTPError } from 'ky';
	import { isHttpError } from '@sveltejs/kit';

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

		let post = ky
			.post('/api/users/', {
				json: {
					username: signupValue.username,
					password: signupValue.password
				}
			})
			.json();

		toast.promise(post, {
			loading: 'Signing up...',
			success: 'Signed up successfully!',
			error: (error) => {
				if (
					isHTTPError(error) &&
					typeof error.data === 'object' &&
					error.data !== null &&
					'message' in error.data
				) {
					return error.data.message;
				} else {
					return 'An error occurred';
				}
			}
		});
	}
</script>

<H1>Authenticate</H1>
<div class="grid grid-cols-2 gap-16">
	<form class={formClasses} onsubmit={signup}>
		<H2>Sign Up</H2>
		<TextInput label="Username" bind:value={signupValue.username} />
		<TextInput label="Password" bind:value={signupValue.password} />
		<TextInput label="Verify Password" bind:value={signupValue.verifyPassword} />
		<Submit value="Register" />
	</form>
	<form class={formClasses}>
		<H2>Login</H2>
		<TextInput label="Username" bind:value={loginValue.username} />
		<TextInput label="Password" bind:value={loginValue.password} />
		<Submit value="Register" />
	</form>
</div>
