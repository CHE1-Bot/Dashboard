<script>
  import { onMount, afterUpdate } from 'svelte';
  import { sidebarLinks, routes, getRoute } from './data/pages.js';

  let hash = window.location.hash || '#/dashboard';
  let currentRoute = getRoute(hash);
  let route = routes[currentRoute] || routes.dashboard;
  let Component;

  const defaultTranscripts = [
    { id: '001', content: 'Sample transcript content for ticket 001.' },
    { id: '002', content: 'Sample transcript content for ticket 002.' },
  ];

  function updateRoute() {
    hash = window.location.hash || '#/dashboard';
    currentRoute = getRoute(hash);
    route = routes[currentRoute] || routes.dashboard;
    loadComponent();
    initializePage();
  }

  function loadComponent() {
    if (route.componentName) {
      import(`./components/${route.componentName}.svelte`).then(module => {
        Component = module.default;
      }).catch(err => {
        console.error('Failed to load component', err);
        Component = null;
      });
    } else {
      Component = null;
    }
  }

  function isSidebarActive(id) {
    return route.sidebarActive === id;
  }

  function isRightNavActive(path) {
    return currentRoute === path;
  }

  function displayTranscripts(transcripts) {
    const transcriptsList = document.getElementById('transcripts-list');
    if (!transcriptsList) return;

    if (transcripts && transcripts.length > 0) {
      transcriptsList.innerHTML = transcripts
        .map(
          (t) => `
            <div class="transcript-item">
              <h3>Ticket ID: ${t.id}</h3>
              <div class="transcript-content">${t.content}</div>
            </div>`
        )
        .join('');
    } else {
      transcriptsList.innerHTML = '<p>No transcripts found.</p>';
    }
  }

  function sortTranscripts(sortType) {
    loadAllTranscripts();
  }

  function searchTranscript(id) {
    const transcriptsList = document.getElementById('transcripts-list');
    if (!transcriptsList) return;

    fetch(`/api/transcripts/${id}`)
      .then((response) => {
        if (response.ok) return response.json();
        throw new Error('Not found');
      })
      .then((data) => displayTranscripts([data]))
      .catch(() => displayTranscripts([]));
  }

  function loadAllTranscripts() {
    const transcriptsList = document.getElementById('transcripts-list');
    if (!transcriptsList) return;

    fetch('/api/transcripts')
      .then((response) => response.json())
      .then((data) => displayTranscripts(data))
      .catch(() => displayTranscripts(defaultTranscripts));
  }

  function initializePage() {
    setTimeout(() => {
      document
        .querySelectorAll('.accordion-header1, .accordion-header2, .accordion-header3, .accordion-header4')
        .forEach((header) => {
          header.onclick = () => header.parentElement.classList.toggle('active');
          header.style.cursor = 'pointer';
          if (!header.hasAttribute('role')) {
            header.setAttribute('role', 'button');
          }
        });

      document.querySelectorAll('.dropdown-button').forEach((button) => {
        button.onclick = () => {
          button.parentElement.classList.toggle('active');
        };
      });

      document.querySelectorAll('.dropdown-content').forEach((content) => {
        content.querySelectorAll('a').forEach((link) => {
          link.onclick = (event) => {
            event.preventDefault();
            const sortType = link.getAttribute('data-sort');
            if (sortType) sortTranscripts(sortType);
            content.parentElement.classList.remove('active');
          };
        });
      });

      const searchBtn = document.getElementById('search-btn');
      const ticketIdInput = document.getElementById('ticket-id');
      if (searchBtn && ticketIdInput) {
        searchBtn.onclick = () => {
          const ticketId = ticketIdInput.value.trim();
          if (ticketId) {
            searchTranscript(ticketId);
          } else {
            loadAllTranscripts();
          }
        };

        ticketIdInput.onkeypress = (event) => {
          if (event.key === 'Enter') {
            searchBtn.click();
          }
        };

        loadAllTranscripts();
      }
    }, 0);
  }

  onMount(() => {
    window.addEventListener('hashchange', updateRoute);
    initializePage();
    loadComponent();
  });

  afterUpdate(() => {
    initializePage();
  });
</script>

<div class="dashboard">
  <aside class="sidebar">
    <div class="logo">
      <img src="/images/Che1logo.png" class="logo-img" alt="Che1 logo" />
      <span>Che1</span>
    </div>

    <nav>
      {#each sidebarLinks as link}
        <a href={`#/${link.path}`} class:selected={isSidebarActive(link.id)}>
          <i class={`fa-solid ${link.icon}`}></i> {link.label}
        </a>
      {/each}
    </nav>
  </aside>

  <div class="main">
    <header class="topbar">
      <h2>{route.pageTitle}</h2>
      <div class="user">
        <img src="https://cdn.discordapp.com/embed/avatars/0.png" alt="User avatar" />
      </div>
    </header>

    <div class="content">
      {#if Component}
        <svelte:component this={Component} />
      {:else}
        {@html route.contentHtml}
      {/if}
    </div>
  </div>

  {#if route.rightNav && route.rightNav.length}
    <aside class="right-sidebar">
      <nav>
        {#each route.rightNav as item}
          <a href={`#/${item.path}`} class:selected={isRightNavActive(item.path)}>
            <i class={`fa-solid ${item.icon}`}></i> {item.label}
          </a>
        {/each}
      </nav>
    </aside>
  {/if}
</div>
