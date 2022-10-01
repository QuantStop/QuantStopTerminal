module.exports = {
    branches: [
        'release',
        'next',
        'next-major',
        '+([0-9])?(.{+([0-9]),x}).x',
        { name: 'beta', prerelease: true },
        { name: 'alpha', prerelease: true },
    ],
    plugins: [
        '@semantic-release/commit-analyzer',
        '@semantic-release/release-notes-generator',
        [
            '@semantic-release/changelog',
            {
                changelogFile: './docs/CHANGELOG.md'
            }
        ],
        '@semantic-release/github',
        [
            '@semantic-release/git',
            {
                assets: ['./docs/CHANGELOG.md', "./package.json", "./package-lock.json"],
                message: 'chore(release): set `package.json` to ${nextRelease.version}\n\n${nextRelease.notes}'
            }
        ]
    ]
}
