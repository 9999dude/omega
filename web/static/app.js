(() => {
  const sidebar = document.querySelector("[data-sidebar]");
  const closeSidebar = () => sidebar?.classList.remove("open");
  document.querySelector("[data-menu-toggle]")?.addEventListener("click", () => sidebar?.classList.toggle("open"));
  document.querySelector("[data-sidebar-scrim]")?.addEventListener("click", closeSidebar);
  sidebar?.querySelectorAll("a").forEach((link) => link.addEventListener("click", closeSidebar));

  document.querySelectorAll(".markdown-body table").forEach((table) => {
    const wrapper = document.createElement("div");
    wrapper.className = "table-wrapper";
    table.parentNode.insertBefore(wrapper, table);
    wrapper.appendChild(table);
  });

  const diagramDialog = document.querySelector("[data-diagram-dialog]");
  const diagramCanvas = document.querySelector("[data-diagram-canvas]");
  const diagramViewport = document.querySelector("[data-diagram-viewport]");
  const diagramZoomOutput = document.querySelector("[data-diagram-zoom]");
  let diagramZoom = 1;
  let diagramTrigger = null;

  const setDiagramZoom = (nextZoom) => {
    diagramZoom = Math.min(3, Math.max(.5, nextZoom));
    const svg = diagramCanvas?.querySelector("svg");
    if (svg) svg.style.width = `${diagramZoom * 100}%`;
    if (diagramZoomOutput) diagramZoomOutput.value = `${Math.round(diagramZoom * 100)}%`;
  };
  const closeDiagram = () => {
    if (!diagramDialog || diagramDialog.hidden) return;
    diagramDialog.hidden = true;
    document.body.classList.remove("diagram-open");
    diagramCanvas?.replaceChildren();
    diagramTrigger?.focus();
  };
  const openDiagram = (stage, trigger) => {
    const svg = stage.querySelector("svg");
    if (!svg || !diagramDialog || !diagramCanvas) return;
    diagramTrigger = trigger || stage;
    const clone = svg.cloneNode(true);
    clone.removeAttribute("width");
    clone.removeAttribute("height");
    clone.style.maxWidth = "none";
    diagramCanvas.replaceChildren(clone);
    diagramDialog.hidden = false;
    document.body.classList.add("diagram-open");
    setDiagramZoom(1);
    diagramViewport?.scrollTo({ top: 0, left: 0 });
    diagramDialog.querySelector("[data-diagram-zoom-in]")?.focus();
  };

  document.querySelectorAll("[data-diagram-close]").forEach((button) => button.addEventListener("click", closeDiagram));
  document.querySelector("[data-diagram-zoom-in]")?.addEventListener("click", () => setDiagramZoom(diagramZoom + .25));
  document.querySelector("[data-diagram-zoom-out]")?.addEventListener("click", () => setDiagramZoom(diagramZoom - .25));
  document.querySelector("[data-diagram-reset]")?.addEventListener("click", () => {
    setDiagramZoom(1);
    diagramViewport?.scrollTo({ top: 0, left: 0, behavior: "smooth" });
  });

  const renderMermaidDiagrams = async () => {
    const blocks = [...document.querySelectorAll(".markdown-body pre > code.language-mermaid")];
    if (!blocks.length || !globalThis.mermaid) return;

    globalThis.mermaid.initialize({
      startOnLoad: false,
      securityLevel: "strict",
      suppressErrorRendering: true,
      theme: "base",
      fontFamily: "Inter, ui-sans-serif, system-ui, sans-serif",
      themeVariables: {
        background: "#32302f",
        primaryColor: "#504945",
        primaryTextColor: "#ebdbb2",
        primaryBorderColor: "#d79921",
        secondaryColor: "#3c3836",
        secondaryTextColor: "#ebdbb2",
        secondaryBorderColor: "#689d6a",
        tertiaryColor: "#665c54",
        tertiaryTextColor: "#fbf1c7",
        tertiaryBorderColor: "#b16286",
        lineColor: "#a89984",
        textColor: "#ebdbb2",
        mainBkg: "#504945",
        nodeBorder: "#d79921",
        clusterBkg: "#3c3836",
        clusterBorder: "#665c54",
        titleColor: "#fbf1c7",
        edgeLabelBackground: "#32302f",
        actorBkg: "#504945",
        actorBorder: "#d79921",
        actorTextColor: "#ebdbb2",
        actorLineColor: "#7c6f64",
        signalColor: "#ebdbb2",
        signalTextColor: "#ebdbb2",
        labelBoxBkgColor: "#3c3836",
        labelBoxBorderColor: "#665c54",
        labelTextColor: "#ebdbb2",
        loopTextColor: "#ebdbb2",
        noteBkgColor: "#504945",
        noteBorderColor: "#d79921",
        noteTextColor: "#ebdbb2",
        activationBkgColor: "#665c54",
        activationBorderColor: "#d79921"
      },
      flowchart: { curve: "basis", htmlLabels: true },
      sequence: { useMaxWidth: true, wrap: true }
    });

    for (const [index, code] of blocks.entries()) {
      const source = code.textContent.trim();
      const original = code.parentElement;
      const card = document.createElement("section");
      card.className = "mermaid-card";
      card.setAttribute("aria-label", "Mermaid diagram");

      const header = document.createElement("header");
      header.className = "mermaid-card-header";
      const label = document.createElement("span");
      label.textContent = "Interactive diagram";
      const openButton = document.createElement("button");
      openButton.type = "button";
      openButton.className = "diagram-open-button";
      openButton.innerHTML = '<span aria-hidden="true">⛶</span> Open and zoom';
      header.append(label, openButton);

      const stage = document.createElement("div");
      stage.className = "mermaid-stage";
      stage.tabIndex = 0;
      stage.setAttribute("role", "button");
      stage.setAttribute("aria-label", "Open diagram viewer");
      card.append(header, stage);

      try {
        const id = `omega-mermaid-${Date.now()}-${index}`;
        const { svg, bindFunctions } = await globalThis.mermaid.render(id, source);
        stage.innerHTML = svg;
        bindFunctions?.(stage);
        original.replaceWith(card);
        openButton.addEventListener("click", () => openDiagram(stage, openButton));
        stage.addEventListener("click", () => openDiagram(stage, stage));
        stage.addEventListener("keydown", (event) => {
          if (event.key === "Enter" || event.key === " ") {
            event.preventDefault();
            openDiagram(stage, stage);
          }
        });
      } catch (error) {
        card.classList.add("mermaid-error");
        card.textContent = `Unable to render this diagram.\n\n${error.message || error}`;
        original.replaceWith(card);
      }
    }
  };
  void renderMermaidDiagrams();

  document.querySelectorAll(".markdown-body pre").forEach((block) => {
    if (block.querySelector("code.language-mermaid")) return;
    const code = block.querySelector("code");
    const languageClass = [...(code?.classList || [])].find((name) => name.startsWith("language-"));
    const sourceLanguage = languageClass?.slice("language-".length).toLowerCase() || "";
    const language = sourceLanguage === "yml" ? "yaml" : sourceLanguage;
    const highlightedLanguages = new Set(["go", "yaml", "json"]);
    const languageLabels = { go: "Go", yaml: "YAML", json: "JSON" };

    if (code && highlightedLanguages.has(language) && globalThis.hljs?.getLanguage(language)) {
      const result = globalThis.hljs.highlight(code.textContent, { language, ignoreIllegals: true });
      code.innerHTML = result.value;
      code.classList.add("hljs");
      block.classList.add("highlighted-code");

      const badge = document.createElement("span");
      badge.className = "code-language";
      badge.textContent = languageLabels[language];
      block.appendChild(badge);
    }

    const button = document.createElement("button");
    button.type = "button";
    button.className = "copy-code";
    button.textContent = "Copy";
    button.addEventListener("click", async () => {
      await navigator.clipboard.writeText(code?.innerText || block.innerText);
      button.textContent = "Copied";
      window.setTimeout(() => { button.textContent = "Copy"; }, 1400);
    });
    block.appendChild(button);
  });

  const progress = document.querySelector(".reading-progress span");
  if (progress && document.body.classList.contains("article-page")) {
    const updateProgress = () => {
      const max = document.documentElement.scrollHeight - window.innerHeight;
      progress.style.width = `${max > 0 ? Math.min(100, window.scrollY / max * 100) : 0}%`;
    };
    window.addEventListener("scroll", updateProgress, { passive: true });
    updateProgress();
  }

  const tocLinks = [...document.querySelectorAll(".toc nav a")];
  if (tocLinks.length) {
    const headings = tocLinks.map((link) => document.getElementById(decodeURIComponent(link.getAttribute("href").slice(1)))).filter(Boolean);
    const observer = new IntersectionObserver((entries) => {
      const visible = entries.filter((entry) => entry.isIntersecting).sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top)[0];
      if (!visible) return;
      tocLinks.forEach((link) => link.classList.toggle("active", link.getAttribute("href") === `#${visible.target.id}`));
    }, { rootMargin: "-90px 0px -72% 0px" });
    headings.forEach((heading) => observer.observe(heading));
  }

  const dialog = document.querySelector("[data-search-dialog]");
  const input = document.querySelector("[data-search-input]");
  const results = document.querySelector("[data-search-results]");
  let activeIndex = -1;
  let requestNumber = 0;

  const openSearch = () => {
    if (!dialog) return;
    dialog.hidden = false;
    document.body.classList.add("search-open");
    window.setTimeout(() => input?.focus(), 0);
  };
  const closeSearch = () => {
    if (!dialog) return;
    dialog.hidden = true;
    document.body.classList.remove("search-open");
    activeIndex = -1;
  };
  const setActive = (index) => {
    const links = [...results.querySelectorAll(".search-result")];
    if (!links.length) return;
    activeIndex = (index + links.length) % links.length;
    links.forEach((link, i) => link.classList.toggle("active", i === activeIndex));
    links[activeIndex].scrollIntoView({ block: "nearest" });
  };
  const escapeHTML = (value) => {
    const node = document.createElement("span");
    node.textContent = value;
    return node.innerHTML;
  };
  const renderResults = (items) => {
    activeIndex = -1;
    if (!items.length) {
      results.innerHTML = '<div class="search-empty"><span>Ω</span><p>No matching notes found.</p></div>';
      return;
    }
    results.innerHTML = items.map((item) => `<a class="search-result" href="${encodeURI(item.url)}"><small>${escapeHTML(item.section)}</small><strong>${escapeHTML(item.title)}</strong><span>${escapeHTML(item.description)}</span></a>`).join("");
  };

  document.querySelectorAll("[data-search-open]").forEach((button) => button.addEventListener("click", openSearch));
  document.querySelectorAll("[data-search-close]").forEach((button) => button.addEventListener("click", closeSearch));
  input?.addEventListener("input", async () => {
    const query = input.value.trim();
    const currentRequest = ++requestNumber;
    if (!query) {
      results.innerHTML = '<div class="search-empty"><span>Ω</span><p>Start typing to explore the library.</p></div>';
      return;
    }
    const response = await fetch(`/api/search?q=${encodeURIComponent(query)}`);
    const items = await response.json();
    if (currentRequest === requestNumber) renderResults(items);
  });
  document.addEventListener("keydown", (event) => {
    if (diagramDialog && !diagramDialog.hidden) {
      if (event.key === "Escape") closeDiagram();
      if ((event.metaKey || event.ctrlKey) && (event.key === "+" || event.key === "=")) setDiagramZoom(diagramZoom + .25);
      if ((event.metaKey || event.ctrlKey) && event.key === "-") setDiagramZoom(diagramZoom - .25);
      return;
    }
    if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "k") {
      event.preventDefault();
      dialog?.hidden ? openSearch() : closeSearch();
      return;
    }
    if (dialog?.hidden) return;
    if (event.key === "Escape") closeSearch();
    if (event.key === "ArrowDown") { event.preventDefault(); setActive(activeIndex + 1); }
    if (event.key === "ArrowUp") { event.preventDefault(); setActive(activeIndex - 1); }
    if (event.key === "Enter" && activeIndex >= 0) results.querySelectorAll(".search-result")[activeIndex]?.click();
  });
})();
