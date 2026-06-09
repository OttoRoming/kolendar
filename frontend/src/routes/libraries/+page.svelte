<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import ky, { isHTTPError } from 'ky';

	ky('/api/libraries')
		.json()
		.then((response) => {
			console.log(response);
		})
		.catch((error) => {
			if (isHTTPError(error)) {
				if (error.response.status == 401) {
					goto(resolve('/login'));
				}
			}
		});
</script>

<div></div>
