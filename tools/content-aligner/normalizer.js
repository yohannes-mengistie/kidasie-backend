(function (root, factory) {
  "use strict";

  const api = factory();

  if (typeof module === "object" && module.exports) {
    module.exports = api;
  }

  root.KidasieNormalizer = api;
})(typeof globalThis !== "undefined" ? globalThis : this, function () {
  "use strict";

  const nonSpokenKinds = new Set([
    "header",
    "anaphora_header",
    "instruction"
  ]);

  const readerKinds = new Set([
    "reading_announcement",
    "reading_ref",
    "scripture"
  ]);

  function clean(value) {
    return typeof value === "string" ? value.trim() : "";
  }

  function joinUnique(values) {
    const seen = new Set();

    return values
      .map(clean)
      .filter(function (value) {
        if (!value || seen.has(value)) {
          return false;
        }

        seen.add(value);
        return true;
      })
      .join("\n\n");
  }

  function normalizeRole(role, kind, allowedRoles) {
    if (role === "people") {
      return "congregation";
    }

    if (role && allowedRoles.has(role)) {
      return role;
    }

    if (nonSpokenKinds.has(kind)) {
      return "rubric";
    }

    if (readerKinds.has(kind)) {
      return "reader";
    }

    return "";
  }

  function baseWarning(page, detail) {
    const warnings = [];

    if (detail) {
      warnings.push(detail);
    }

    if (clean(page.note)) {
      warnings.push("Source note: " + clean(page.note));
    }

    if (clean(page.instruction)) {
      warnings.push("Instruction: " + clean(page.instruction));
    }

    if (clean(page.instruction_amharic)) {
      warnings.push("Instruction: " + clean(page.instruction_amharic));
    }

    if (page.needs_review !== false) {
      warnings.push("Source marks this page as needing review.");
    }

    return warnings.join(" ");
  }

  function missingLanguageWarning(textGeez, textAmharic) {
    const missing = [];

    if (!clean(textGeez)) {
      missing.push("Geʽez");
    }

    if (!clean(textAmharic)) {
      missing.push("Amharic");
    }

    if (missing.length === 0) {
      return "";
    }

    return "Add reviewed " + missing.join(" and ") +
      " text before final export.";
  }

  function makeSegment(page, values) {
    const languageWarning = missingLanguageWarning(
      values.text_geez,
      values.text_am
    );
    const warning = [baseWarning(page, values.detail), languageWarning]
      .filter(Boolean)
      .join(" ");

    return {
      source_page: page.number,
      source_part: values.source_part || "",
      kind: values.kind || page.kind || "content",
      warning: warning || "Confirm the corrected text and speaker against audio.",
      include: values.include !== false,
      role: values.role || "",
      text_geez: clean(values.text_geez),
      text_am: clean(values.text_am),
      text_en: clean(values.text_en),
      start_ms: null,
      end_ms: null
    };
  }

  function hasText(values) {
    return Boolean(
      clean(values.text_geez) ||
      clean(values.text_am) ||
      clean(values.text_en)
    );
  }

  function mainText(page) {
    return {
      text_geez: page.text_geez,
      text_am: page.text_amharic,
      text_en: page.text_english
    };
  }

  function titleText(page) {
    return {
      text_geez: page.title_geez,
      text_am: page.title_amharic,
      text_en: page.title_english || page.title
    };
  }

  function responseText(page) {
    return {
      text_geez: joinUnique([
        page.response_people_geez,
        page.response_people,
        page.text_geez_people
      ]),
      text_am: joinUnique([
        page.response_people_amharic,
        page.response_amharic,
        page.text_amharic_people
      ]),
      text_en: joinUnique([
        page.response_people_english,
        page.response_english,
        page.text_english_people
      ])
    };
  }

  function appendPart(segments, page, part, index, allowedRoles, label) {
    const values = {
      text_geez: part.text_geez,
      text_am: part.text_amharic,
      text_en: part.text_english
    };

    if (!hasText(values)) {
      return;
    }

    const instructionReference =
      page.kind === "instruction" && !clean(values.text_geez);

    segments.push(makeSegment(page, {
      source_part: label + "-" + (index + 1),
      kind: page.kind + ":" + label,
      detail: "Split from " + label + " " + (index + 1) +
        " on the corrected slide." +
        (instructionReference
          ? " Retained as a non-spoken instruction."
          : ""),
      include: !instructionReference,
      role: instructionReference
        ? "rubric"
        : normalizeRole(part.role, page.kind, allowedRoles),
      text_geez: values.text_geez,
      text_am: values.text_am,
      text_en: values.text_en
    }));
  }

  function curatedPageToSegments(page, allowedRoles) {
    const segments = [];
    const main = mainText(page);
    const titles = titleText(page);
    const hasMain = hasText(main);
    const hasTitle = hasText(titles);
    const role = normalizeRole(page.role, page.kind, allowedRoles);

    if (hasMain) {
      segments.push(makeSegment(page, {
        source_part: "main",
        kind: page.kind,
        detail: page.role === "mixed"
          ? "This main block has mixed speakers; confirm or split it."
          : "",
        include: !nonSpokenKinds.has(page.kind),
        role: role,
        text_geez: main.text_geez,
        text_am: main.text_am,
        text_en: main.text_en
      }));
    } else if (hasTitle) {
      segments.push(makeSegment(page, {
        source_part: "title",
        kind: page.kind + ":title",
        detail: "Heading retained for reference and excluded from audio export.",
        include: false,
        role: "rubric",
        text_geez: titles.text_geez,
        text_am: titles.text_am,
        text_en: titles.text_en
      }));
    }

    (Array.isArray(page.parts) ? page.parts : []).forEach(
      function (part, index) {
        appendPart(segments, page, part, index, allowedRoles, "part");
      }
    );

    const deaconInstruction = {
      text_geez: page.deacon_instruction_geez,
      text_am: page.deacon_instruction_amharic,
      text_en: page.deacon_instruction_english
    };

    if (hasText(deaconInstruction)) {
      segments.push(makeSegment(page, {
        source_part: "deacon-instruction",
        kind: page.kind + ":deacon-instruction",
        detail: "Separate deacon instruction found on the corrected slide.",
        include: true,
        role: "deacon",
        text_geez: deaconInstruction.text_geez,
        text_am: deaconInstruction.text_am,
        text_en: deaconInstruction.text_en
      }));
    }

    const response = responseText(page);

    if (hasText(response)) {
      segments.push(makeSegment(page, {
        source_part: "people-response",
        kind: page.kind + ":people-response",
        detail: "Separate people response found on the corrected slide.",
        include: true,
        role: "congregation",
        text_geez: response.text_geez,
        text_am: response.text_am,
        text_en: response.text_en
      }));
    }

    (Array.isArray(page.response_mixed) ? page.response_mixed : []).forEach(
      function (part, index) {
        appendPart(
          segments,
          page,
          part,
          index,
          allowedRoles,
          "mixed-response"
        );
      }
    );

    return segments;
  }

  function isCuratedPage(page) {
    return [
      "text_geez",
      "text_amharic",
      "text_english",
      "title_geez",
      "title_amharic",
      "title_english",
      "parts",
      "response_people",
      "response_mixed"
    ].some(function (key) {
      return Object.prototype.hasOwnProperty.call(page, key);
    });
  }

  function legacyPagesToSegments(draft, ocrLanguage) {
    return (draft.pages || [])
      .filter(function (page) {
        if (page.kind === "role_header" || page.kind === "empty") {
          return false;
        }

        return Boolean(page.text_ethiopic_ocr || page.text_english_ocr);
      })
      .map(function (page) {
        return {
          source_page: page.number,
          source_part: "ocr",
          kind: page.kind,
          warning: page.extraction_warning || "",
          include: true,
          role: page.role_suggestion || "",
          text_geez:
            ocrLanguage === "geez" ? page.text_ethiopic_ocr || "" : "",
          text_am:
            ocrLanguage === "amharic" ? page.text_ethiopic_ocr || "" : "",
          text_en: page.text_english_ocr || "",
          start_ms: null,
          end_ms: null
        };
      });
  }

  function pagesToSegments(draft, ocrLanguage, allowedRoles) {
    const pages = Array.isArray(draft.pages) ? draft.pages : [];

    if (!pages.some(isCuratedPage)) {
      return legacyPagesToSegments(draft, ocrLanguage);
    }

    return pages.flatMap(function (page) {
      return curatedPageToSegments(page, allowedRoles);
    });
  }

  return {
    pagesToSegments: pagesToSegments
  };
});
