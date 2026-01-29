
import fs from 'fs';
import path from 'path';
import https from 'https';

const INBOX_PATH = path.join(process.cwd(), 'docs', 'INBOX_LINKS.md');

async function main() {
    console.log("📂 Reading INBOX_LINKS.md...");
    if (!fs.existsSync(INBOX_PATH)) {
        console.error("❌ File not found:", INBOX_PATH);
        return;
    }

    const content = fs.readFileSync(INBOX_PATH, 'utf-8');
    const lines = content.split('\n');

    const links: { line: number, url: string, checked: boolean, domain: string }[] = [];
    const urlRegex = /(https?:\/\/[^\s\)]+)/g;

    lines.forEach((line, index) => {
        const match = line.match(urlRegex);
        if (match) {
            const url = match[0];
            const checked = line.trim().startsWith('- [x]');
            try {
                const domain = new URL(url).hostname;
                links.push({ line: index + 1, url, checked, domain });
            } catch (e) {
                // Invalid URL
            }
        }
    });

    let report = "# Inbox Triage Report\n\n";
    report += `🔍 Found ${links.length} total links.\n`;
    const unchecked = links.filter(l => !l.checked);
    report += `📝 ${unchecked.length} unchecked items.\n\n`;

    // Group by Domain
    const domainCounts: Record<string, number> = {};
    unchecked.forEach(l => {
        const domain = l.domain || 'unknown';
        domainCounts[domain] = (domainCounts[domain] || 0) + 1;
    });

    report += "## Top Domains (Unchecked)\n";
    Object.entries(domainCounts)
        .sort((a, b) => b[1] - a[1])
        .slice(0, 15)
        .forEach(([domain, count]) => report += `- **${domain}**: ${count}\n`);

    report += "\n## GitHub Repositories (Potential Submodules)\n";
    const githubLinks = unchecked.filter(l => l.domain === 'github.com');
    githubLinks.forEach(l => {
        // Extract user/repo
        const parts = l.url.split('github.com/')[1]?.split('/');
        if (parts && parts.length >= 2) {
            report += `- [ ] [${parts[0]}/${parts[1]}](${l.url})\n`;
        } else {
            report += `- [ ] ${l.url}\n`;
        }
    });

    report += "\n## Other Links\n";
    unchecked.filter(l => l.domain !== 'github.com').forEach(l => {
        report += `- [ ] ${l.url}\n`;
    });

    const reportPath = path.join(process.cwd(), 'docs', 'TRIAGE_REPORT.md');
    fs.writeFileSync(reportPath, report);
    console.log(`✅ Report generated at ${reportPath}`);
}

main().catch(console.error);
