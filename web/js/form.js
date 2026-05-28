const MAX_TEXT = 6000;
const MAX_OTHER = 100;

const STAMP_SPROUT_ICON = `<svg class="btn-stamp__glyph" viewBox="0 0 16 16" width="16" height="16" aria-hidden="true" focusable="false"><path d="M8 13.25V7.75" stroke="currentColor" stroke-width="1.2" stroke-linecap="round"/><path d="M8 8.25C8 8.25 5.4 7.5 4.25 4.75C6.35 5.85 8 8.25 8 8.25Z" fill="currentColor"/><path d="M8 8.25C8 8.25 10.6 7.5 11.75 4.75C9.65 5.85 8 8.25 8 8.25Z" fill="currentColor"/><ellipse cx="8" cy="13.35" rx="0.95" ry="0.62" fill="currentColor"/></svg>`;

const GROUPS = [
  { name: "types", otherId: "typesOther", trayId: "tray-type" },
  { name: "genres", otherId: "genresOther", trayId: "tray-genre" },
  { name: "businessModels", otherId: "businessModelsOther", trayId: "tray-model" },
  { name: "barriers", otherId: "barriersOther", trayId: "tray-barriers" },
];

const SECTIONS = [
  "tray-type",
  "tray-genre",
  "tray-model",
  "tray-barriers",
  "tray-idea",
  "tray-mvp",
];

const TRAY_META = {
  "tray-type": { group: "types", kind: "checkbox", title: "Type" },
  "tray-genre": { group: "genres", kind: "checkbox", title: "Genre" },
  "tray-model": { group: "businessModels", kind: "checkbox", title: "Business model" },
  "tray-barriers": { group: "barriers", kind: "checkbox", title: "Blockers" },
  "tray-idea": { field: "idea", kind: "text", title: "Idea" },
  "tray-mvp": { field: "mvp", kind: "text", title: "MVP" },
};

const visitedSteps = new Set();
const touchedSteps = new Set();
let currentStepIndex = 0;
let prefersReduced = false;
let approvalFadeTimer = null;

function $(id) {
  return document.getElementById(id);
}

function selectedValues(groupName) {
  return Array.from(
    document.querySelectorAll(`input[name="${groupName}"]:checked`)
  ).map((el) => el.value);
}

function hasOther(groupName) {
  return selectedValues(groupName).includes("Other");
}

function validateSection(trayId) {
  const meta = TRAY_META[trayId];
  if (!meta) {
    return { valid: false, message: "Unknown section." };
  }

  if (meta.kind === "checkbox") {
    const vals = selectedValues(meta.group);
    if (vals.length === 0) {
      return { valid: false, message: "Select at least one option." };
    }
    if (vals.includes("Other")) {
      const group = GROUPS.find((g) => g.name === meta.group);
      const other = $(group.otherId).value.trim();
      if (!other) {
        return { valid: false, message: "Describe your Other choice." };
      }
      if (other.length > MAX_OTHER) {
        return { valid: false, message: "Other description is too long." };
      }
    }
    return { valid: true };
  }

  const value = $(meta.field).value.trim();
  if (!value) {
    return { valid: false, message: "This section is required." };
  }
  if (value.length > MAX_TEXT) {
    return { valid: false, message: "Text exceeds the maximum length." };
  }
  return { valid: true };
}

function sectionStatus(trayId) {
  if (validateSection(trayId).valid) {
    return { label: "Complete", className: "status-complete" };
  }
  return { label: "Needs answer", className: "status-needs" };
}

function allSectionsComplete() {
  return SECTIONS.every((id) => validateSection(id).valid);
}

function updateStatusUI() {
  const currentTrayId = SECTIONS[currentStepIndex];

  SECTIONS.forEach((trayId) => {
    const status = sectionStatus(trayId);
    const badge = document.querySelector(`[data-step-badge="${trayId}"]`);
    const showInlineStatus =
      trayId === currentTrayId && status.className === "status-needs";

    if (badge) {
      badge.textContent = status.label;
      badge.className = `step-badge ${status.className}`;
      badge.hidden = !showInlineStatus;
    }

    document.querySelectorAll(`[data-step-status="${trayId}"]`).forEach((el) => {
      el.textContent = status.label;
      el.className = `step-status ${status.className}`;
    });

    document.querySelectorAll(`a[href="#${trayId}"]`).forEach((link) => {
      link.classList.remove("step-complete", "step-needs", "step-progress", "step-required");
      link.classList.add(`step-${status.className.replace("status-", "")}`);
    });

    const step = document.getElementById(trayId);
    if (step) {
      step.classList.toggle("is-complete", status.className === "status-complete");
      step.classList.toggle("is-needs-answer", status.className === "status-needs");
    }
  });

  const progressLabel = $("stepProgressLabel");
  if (progressLabel) {
    progressLabel.textContent = `Step ${currentStepIndex + 1} of ${SECTIONS.length}`;
  }

  const submitBtn = $("submitBtn");
  const submitHint = $("submitHint");
  const finale = $("formFinale");
  const complete = allSectionsComplete();

  if (submitBtn) {
    submitBtn.disabled = !complete;
  }
  if (finale) {
    finale.hidden = !complete;
  }
  if (submitHint) {
    submitHint.hidden = complete;
  }
}

function setActiveNav(trayId) {
  document.querySelectorAll(".section-rail a, .step-strip a").forEach((link) => {
    const active = link.getAttribute("href") === `#${trayId}`;
    link.classList.toggle("is-active", active);
    if (active) {
      link.setAttribute("aria-current", "step");
    } else {
      link.removeAttribute("aria-current");
    }
  });
}

function showStepError(trayId, message) {
  const step = document.getElementById(trayId);
  const errorEl = step?.querySelector(".step-error");
  if (!errorEl) return;
  if (message) {
    hideStepApproval(trayId);
    errorEl.textContent = message;
    errorEl.hidden = false;
  } else {
    errorEl.textContent = "";
    errorEl.hidden = true;
  }
}

function hideStepApproval(trayId) {
  if (approvalFadeTimer) {
    clearTimeout(approvalFadeTimer);
    approvalFadeTimer = null;
  }

  const step = document.getElementById(trayId);
  const approvalEl = step?.querySelector(".step-approval");
  if (!approvalEl) return;

  approvalEl.textContent = "";
  approvalEl.hidden = true;
  approvalEl.classList.remove("is-visible", "is-fading");
}

function showStepApproval(trayId, message, onComplete) {
  const step = document.getElementById(trayId);
  const approvalEl = step?.querySelector(".step-approval");
  if (!approvalEl) {
    onComplete?.();
    return;
  }

  if (approvalFadeTimer) {
    clearTimeout(approvalFadeTimer);
    approvalFadeTimer = null;
  }

  approvalEl.textContent = message;
  approvalEl.hidden = false;
  approvalEl.classList.remove("is-fading");
  approvalEl.classList.add("is-visible");

  const holdMs = prefersReduced ? 700 : 1500;
  const fadeMs = prefersReduced ? 150 : 1400;

  approvalFadeTimer = setTimeout(() => {
    approvalEl.classList.add("is-fading");
    approvalFadeTimer = setTimeout(() => {
      hideStepApproval(trayId);
      onComplete?.();
    }, fadeMs);
  }, holdMs);
}

function showStep(index, scroll = true) {
  currentStepIndex = index;
  const trayId = SECTIONS[index];
  visitedSteps.add(trayId);

  SECTIONS.forEach((id, i) => {
    const step = document.getElementById(id);
    if (!step) return;
    step.classList.toggle("is-current", i === index);
  });

  setActiveNav(trayId);
  updateStatusUI();

  const step = document.getElementById(trayId);
  const nav = step?.querySelector(".wizard-nav");
  const prevBtn = step?.querySelector("[data-step-prev]");
  const nextBtn = step?.querySelector("[data-step-next]");
  if (prevBtn) {
    prevBtn.hidden = index === 0;
  }
  if (nav) {
    nav.classList.toggle("wizard-nav--solo", index === 0);
  }
  if (nextBtn) {
    const label = nextBtn.querySelector("[data-step-next-label]");
    if (label) {
      label.textContent =
        index === SECTIONS.length - 1 ? "Review intake" : "Next section";
    }
  }

  if (!scroll) return;
  step?.scrollIntoView({
    behavior: prefersReduced ? "auto" : "smooth",
    block: "start",
  });
}

function goToStep(index, scroll = true) {
  if (index < 0 || index >= SECTIONS.length) return;
  hideStepApproval(SECTIONS[currentStepIndex]);
  showStep(index, scroll);
}

function tryAdvanceStep() {
  const trayId = SECTIONS[currentStepIndex];
  const result = validateSection(trayId);
  if (!result.valid) {
    touchedSteps.add(trayId);
    showStepError(trayId, result.message);
    updateStatusUI();
    return;
  }

  showStepError(trayId, "");
  if (currentStepIndex < SECTIONS.length - 1) {
    goToStep(currentStepIndex + 1, true);
    return;
  }

  updateStatusUI();
  if (allSectionsComplete()) {
    showStepApproval(trayId, "Looking good!", () => {
      $("formFinale")?.scrollIntoView({
        behavior: prefersReduced ? "auto" : "smooth",
        block: "nearest",
      });
    });
    return;
  }

  const incomplete = SECTIONS.find((id) => !validateSection(id).valid);
  if (incomplete) {
    goToStep(SECTIONS.indexOf(incomplete), true);
    showStepError(incomplete, validateSection(incomplete).message);
  }
}

function setupWizardSteps() {
  const form = $("ideaForm");
  form.classList.add("wizard-form");

  SECTIONS.forEach((trayId, index) => {
    const tray = document.getElementById(trayId);
    if (!tray) return;

    tray.classList.add("wizard-step");

    let titleHTML = "";
    const legend = tray.querySelector("legend");
    const heading = tray.querySelector(".tray-heading");
    if (legend) {
      titleHTML = legend.innerHTML;
      legend.classList.add("visually-hidden");
    } else if (heading) {
      titleHTML = heading.innerHTML;
      heading.remove();
    }

    const panel = document.createElement("div");
    panel.className = "step-panel";

    const header = document.createElement("div");
    header.className = "step-header";
    header.innerHTML = `
      <div class="step-heading">
        <h2 class="step-title">${titleHTML}</h2>
        <span class="step-badge status-needs" data-step-badge="${trayId}" hidden>Needs answer</span>
      </div>
    `;

    const error = document.createElement("p");
    error.className = "step-error";
    error.hidden = true;

    const approval = document.createElement("p");
    approval.className = "step-approval";
    approval.hidden = true;
    approval.setAttribute("role", "status");
    approval.setAttribute("aria-live", "polite");

    const nav = document.createElement("div");
    nav.className = `wizard-nav${index === 0 ? " wizard-nav--solo" : ""}`;
    nav.innerHTML = `
      <button type="button" class="btn-secondary btn-nav-prev" data-step-prev hidden>Previous</button>
      <button type="button" class="btn-stamp btn-nav-next" data-step-next>
        <span class="btn-stamp__icon">${STAMP_SPROUT_ICON}</span>
        <span class="btn-stamp__label" data-step-next-label>Next section</span>
      </button>
    `;

    nav.querySelector("[data-step-prev]")?.addEventListener("click", () => {
      goToStep(currentStepIndex - 1, true);
    });
    nav.querySelector("[data-step-next]")?.addEventListener("click", tryAdvanceStep);

    const body = document.createElement("div");
    body.className = "step-body";
    [...tray.children].forEach((child) => {
      if (child === legend) return;
      body.appendChild(child);
    });

    panel.append(header, error, body, approval, nav);
    tray.replaceChildren(panel);
    tray.classList.toggle("is-current", index === 0);
  });

  visitedSteps.add(SECTIONS[0]);
  updateStatusUI();
}

function setupSectionNav() {
  document.querySelectorAll(".section-rail a, .step-strip a").forEach((link) => {
    link.addEventListener("click", (event) => {
      event.preventDefault();
      const trayId = link.getAttribute("href").slice(1);
      const index = SECTIONS.indexOf(trayId);
      if (index >= 0) {
        goToStep(index, true);
      }
    });
  });
}

function markTouched(trayId) {
  touchedSteps.add(trayId);
  updateStatusUI();
}

function setupOtherToggle(groupName, otherInputId) {
  const otherInput = $(otherInputId);
  const otherField = otherInput?.closest(".other-field");
  const group = GROUPS.find((g) => g.name === groupName);

  const sync = () => {
    const show = hasOther(groupName);
    if (otherField) {
      otherField.hidden = !show;
    }
    otherInput.required = show;
    if (!show) {
      otherInput.value = "";
    }
    updateCount(otherInput, MAX_OTHER);
    if (group) markTouched(group.trayId);
  };

  document.querySelectorAll(`input[name="${groupName}"]`).forEach((el) => {
    el.addEventListener("change", sync);
  });
  sync();
}

function setupInteractionListeners() {
  GROUPS.forEach((g) => {
    document.querySelectorAll(`input[name="${g.name}"]`).forEach((input) => {
      input.addEventListener("change", () => {
        markTouched(g.trayId);
        showStepError(g.trayId, "");
      });
    });
    const otherInput = $(g.otherId);
    otherInput?.addEventListener("input", () => {
      markTouched(g.trayId);
      showStepError(g.trayId, "");
    });
  });

  ["idea", "mvp"].forEach((fieldId) => {
    const el = $(fieldId);
    if (!el) return;
    const trayId = Object.entries(TRAY_META).find(([, meta]) => meta.field === fieldId)?.[0];
    el.addEventListener("input", () => {
      if (trayId) {
        markTouched(trayId);
        showStepError(trayId, "");
      }
    });
  });
}

function updateCount(el, max) {
  const counter = document.querySelector(`[data-count-for="${el.id}"]`);
  if (!counter) return;
  counter.textContent = `${el.value.length} / ${max}`;
}

function setupCounters() {
  [
    "idea",
    "mvp",
    "typesOther",
    "genresOther",
    "businessModelsOther",
    "barriersOther",
  ].forEach((id) => {
    const el = $(id);
    if (!el) return;
    const max = id === "idea" || id === "mvp" ? MAX_TEXT : MAX_OTHER;
    el.addEventListener("input", () => updateCount(el, max));
    updateCount(el, max);
  });
}

function validateClient() {
  const errors = [];
  for (const trayId of SECTIONS) {
    const result = validateSection(trayId);
    if (!result.valid) {
      errors.push(`${TRAY_META[trayId].title}: ${result.message}`);
    }
  }
  return errors;
}

function buildPayload() {
  return {
    types: selectedValues("types"),
    typesOther: $("typesOther").value.trim(),
    genres: selectedValues("genres"),
    genresOther: $("genresOther").value.trim(),
    businessModels: selectedValues("businessModels"),
    businessModelsOther: $("businessModelsOther").value.trim(),
    barriers: selectedValues("barriers"),
    barriersOther: $("barriersOther").value.trim(),
    idea: $("idea").value,
    mvp: $("mvp").value,
  };
}

function setStatus(message, kind) {
  const status = $("status");
  status.textContent = message;
  status.className = kind || "";
}

function resetFormState(form) {
  form.reset();
  visitedSteps.clear();
  touchedSteps.clear();

  GROUPS.forEach((g) => {
    const otherInput = $(g.otherId);
    const otherField = otherInput?.closest(".other-field");
    if (otherField) otherField.hidden = true;
    otherInput.required = false;
    otherInput.value = "";
    updateCount(otherInput, MAX_OTHER);
  });

  setupCounters();
  SECTIONS.forEach((trayId) => {
    showStepError(trayId, "");
    hideStepApproval(trayId);
  });
  goToStep(0, false);
  visitedSteps.add(SECTIONS[0]);
  updateStatusUI();
}

async function handleSubmit(event) {
  event.preventDefault();

  if (!allSectionsComplete()) {
    setStatus("Complete all sections before submitting.", "error");
    updateStatusUI();
    return;
  }

  const errors = validateClient();
  if (errors.length) {
    setStatus(errors[0], "error");
    const trayId = SECTIONS.find((id) => !validateSection(id).valid);
    if (trayId) {
      goToStep(SECTIONS.indexOf(trayId), true);
      showStepError(trayId, validateSection(trayId).message);
    }
    return;
  }

  const btn = $("submitBtn");
  btn.disabled = true;
  setStatus("Submitting\u2026", "");

  try {
    const res = await fetch("/api/ideas", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(buildPayload()),
    });
    const data = await res.json().catch(() => ({}));
    if (!res.ok) {
      setStatus(data.error || "Submission failed.", "error");
      return;
    }
    setStatus("Your idea was sent to the garden.", "success");
    resetFormState(event.target);
  } catch {
    setStatus("Network error. Try again.", "error");
  } finally {
    if (btn) {
      btn.disabled = !allSectionsComplete();
    }
  }
}

document.addEventListener("DOMContentLoaded", () => {
  prefersReduced = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
  GROUPS.forEach((g) => setupOtherToggle(g.name, g.otherId));
  setupCounters();
  setupWizardSteps();
  setupSectionNav();
  setupInteractionListeners();
  $("ideaForm").addEventListener("submit", handleSubmit);
});
