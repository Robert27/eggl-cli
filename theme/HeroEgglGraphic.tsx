export function HeroEgglGraphic() {
  return (
    <div className="eggl-hero-graphic">
      <svg
        className="eggl-hero-graphic__svg"
        viewBox="0 0 560 380"
        role="img"
        aria-labelledby="eggl-graphic-title eggl-graphic-description"
      >
        <title id="eggl-graphic-title">eggl terminal session</title>
        <desc id="eggl-graphic-description">
          A simulated terminal session showing eggl CLI checking an environment, switching profiles, and starting a port-forward.
        </desc>
        <rect x="12" y="12" width="536" height="356" rx="16" fill="#17191c" />
        <rect x="12" y="12" width="536" height="44" rx="16" fill="var(--eggl-ink)" />
        <circle cx="38" cy="34" r="5" fill="#ef4444" />
        <circle cx="56" cy="34" r="5" fill="#f59e0b" />
        <circle cx="74" cy="34" r="5" fill="#22c55e" />
        <text x="100" y="39" fill="#c6cbc8" fontSize="13" className="eggl-hero-graphic__mono">eggl CLI / terminal session</text>

        <rect x="28" y="70" width="504" height="264" rx="8" fill="#101214" />
        <path d="M72 82v240" stroke="#ffffff18" />
        <text x="42" y="98" fill="#737b83" fontSize="11" className="eggl-hero-graphic__mono">01</text>
        <text x="82" y="98" fill="var(--eggl-accent-bright)" fontSize="13" className="eggl-hero-graphic__mono">$</text>
        <text x="98" y="98" fill="#f0f1ef" fontSize="13" className="eggl-hero-graphic__mono">eggl doctor</text>

        <text x="42" y="120" fill="#737b83" fontSize="11" className="eggl-hero-graphic__mono">02</text>
        <text x="88" y="120" fill="#22c55e" fontSize="12" className="eggl-hero-graphic__mono">ok</text>
        <text x="116" y="120" fill="#c6cbc8" fontSize="12" className="eggl-hero-graphic__mono">config</text>
        <text x="282" y="120" fill="#22c55e" fontSize="12" className="eggl-hero-graphic__mono">ok</text>
        <text x="310" y="120" fill="#c6cbc8" fontSize="12" className="eggl-hero-graphic__mono">kubectl</text>

        <text x="42" y="140" fill="#737b83" fontSize="11" className="eggl-hero-graphic__mono">03</text>
        <text x="88" y="140" fill="#22c55e" fontSize="12" className="eggl-hero-graphic__mono">ok</text>
        <text x="116" y="140" fill="#c6cbc8" fontSize="12" className="eggl-hero-graphic__mono">git</text>
        <text x="282" y="140" fill="#22c55e" fontSize="12" className="eggl-hero-graphic__mono">ok</text>
        <text x="310" y="140" fill="#c6cbc8" fontSize="12" className="eggl-hero-graphic__mono">tailscale</text>

        <path d="M88 153h422" stroke="#ffffff18" />
        <text x="42" y="176" fill="#737b83" fontSize="11" className="eggl-hero-graphic__mono">04</text>
        <text x="82" y="176" fill="var(--eggl-accent-bright)" fontSize="13" className="eggl-hero-graphic__mono">$</text>
        <text x="98" y="176" fill="#f0f1ef" fontSize="13" className="eggl-hero-graphic__mono">eggl env use homelab</text>

        <text x="42" y="198" fill="#737b83" fontSize="11" className="eggl-hero-graphic__mono">05</text>
        <text x="88" y="198" fill="#737b83" fontSize="12" className="eggl-hero-graphic__mono">profile</text>
        <text x="172" y="198" fill="#c6cbc8" fontSize="12" className="eggl-hero-graphic__mono">work -&gt; homelab</text>
        <text x="42" y="218" fill="#737b83" fontSize="11" className="eggl-hero-graphic__mono">06</text>
        <text x="88" y="218" fill="#737b83" fontSize="12" className="eggl-hero-graphic__mono">kube</text>
        <text x="172" y="218" fill="#c6cbc8" fontSize="12" className="eggl-hero-graphic__mono">context-b -&gt; homelab</text>
        <text x="42" y="238" fill="#737b83" fontSize="11" className="eggl-hero-graphic__mono">07</text>
        <text x="88" y="238" fill="#737b83" fontSize="12" className="eggl-hero-graphic__mono">tailscale</text>
        <text x="172" y="238" fill="#c6cbc8" fontSize="12" className="eggl-hero-graphic__mono">work@example -&gt; homelab</text>

        <path d="M88 251h422" stroke="#ffffff18" />
        <text x="42" y="274" fill="#737b83" fontSize="11" className="eggl-hero-graphic__mono">08</text>
        <text x="82" y="274" fill="var(--eggl-accent-bright)" fontSize="13" className="eggl-hero-graphic__mono">$</text>
        <text x="98" y="274" fill="#f0f1ef" fontSize="13" className="eggl-hero-graphic__mono">eggl pf grafana --open</text>
        <text x="42" y="296" fill="#737b83" fontSize="11" className="eggl-hero-graphic__mono">09</text>
        <text x="88" y="296" fill="#737b83" fontSize="12" className="eggl-hero-graphic__mono">tunnel</text>
        <text x="172" y="296" fill="#c6cbc8" fontSize="12" className="eggl-hero-graphic__mono">monitoring/svc/grafana</text>
        <text x="42" y="316" fill="#737b83" fontSize="11" className="eggl-hero-graphic__mono">10</text>
        <text x="88" y="316" fill="#22c55e" fontSize="12" className="eggl-hero-graphic__mono">ready</text>
        <text x="172" y="316" fill="var(--eggl-accent-bright)" fontSize="12" className="eggl-hero-graphic__mono">http://localhost:3000</text>
        <rect x="344" y="305" width="8" height="15" fill="var(--eggl-accent-bright)" className="eggl-hero-graphic__cursor" />

        <text x="28" y="354" fill="var(--rp-c-text-3)" fontSize="11" className="eggl-hero-graphic__mono">3 commands</text>
        <text x="462" y="354" fill="var(--eggl-success)" fontSize="11" className="eggl-hero-graphic__mono">0 errors</text>
      </svg>
    </div>
  );
}
