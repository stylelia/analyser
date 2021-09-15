package analyser

func HandleEvent() error {
	return handle()
}

func handle() error {
	// Fetch the latest default commit sha and check it against cache

	// Check cache for cookstyle for a given repo.
	// If exists, check version - if equal and if commit sha equal to cache, leave app

	// If not exists or version is different or sha is different, clone the repo

	// run 'cookstyle -a --format json'

	// If cookstyle finds a change, create a new branch 'styleila/cookstyle_<version>'
	// If no change, update cache with cookstyle and default branch sha

	// Raise a PR for that change
	// put in pr body nice message based on json response from cookstyle

	// update cache with default branch sha & cookstyle version

	// see: https://github.com/Xorima/github-cookstyle-runner/blob/main/app/entrypoint.ps1#L139 to L157

	return nil
}
