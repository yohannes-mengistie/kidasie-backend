# Ethiopian Orthodox Kidasie Companion

  Product requirements — draft for review

  Status: Draft 0.1

  Date: 2026-07-29

## 1. Product vision

Build a respectful, modern, and offline-first mobile companion that helps Ethiopian
Orthodox Tewahedo faithful understand, learn, and follow the Divine Liturgy
(Kidasie), its responses, and approved chants.

The product should help a person:

- follow the order of the service without losing their place;
- see approved Geʽez text together with Amharic and, where approved, English;
- know whether a line belongs to the priest, deacon, congregation, or chanter;
- listen and practise with synchronized text and approved audio;
- download content once and continue using it without internet access;
- learn the meaning and pronunciation of the text outside the service.

The app is a companion to Church worship and teaching. It must not present itself
as a substitute for participation in church or as an independent theological
authority.

## 2. Religious and content foundation

The product model must reflect that the Divine Liturgy has a common order and a
variable Anaphora. The Ethiopian Orthodox Tewahedo Church recognizes fourteen
Anaphoras, with the Anaphora of the Apostles commonly used and others associated
with particular occasions. The exact ordering, naming, text, rubrics, and usage
rules in the app must be confirmed by the appointed Church reviewers.

The project must have a named content authority before public release. This can
be a diocese, parish, theological institution, or formally appointed committee of
priests, deacons, and qualified Geʽez scholars.

No liturgical text, translation, explanation, icon, musical recording, or font
may be published until:

1. its source is recorded;
2. the project has permission to use it;
3. a qualified reviewer has approved it;
4. the approved version is immutable and auditable.

Different accepted editions or local practices must not be silently combined.
Where variation is important, the app should identify the edition or tradition
and let the user choose an approved version.

## 3. Target users

### Primary users

- Lay faithful who want to follow and understand Kidasie.
- Young people and diaspora users who need Amharic or English support alongside
  Geʽez.
- Beginners learning congregational responses and pronunciation.

### Secondary users

- Sunday school and theology students.
- Deacons, chanters, and clergy using the app for teaching or preparation.
- Church content reviewers and administrators.

The first release is for learning and lay participation. It is not a replacement
for clergy service books or professional liturgical training.

## 4. Product principles

- **Church-approved:** theological correctness has priority over speed of release.
- **Offline-first:** the core experience works without a connection or account.
- **Low distraction:** worship mode avoids unnecessary controls, notifications,
  advertising, and visual movement.
- **Clear roles:** every spoken or chanted line visibly identifies its role.
- **Multilingual:** translations support understanding without replacing the
  approved source text.
- **Accessible:** Ethiopic text is readable on small and inexpensive phones.
- **Privacy-respecting:** microphone, contacts, and location are not required for
  core use.
- **Versioned:** users can tell which approved content edition they have.

## 5. User experience modes

### 5.1 Learn mode

For study outside the service:

- browse the structure of Kidasie by section;
- display Geʽez, Amharic, and English in one-, two-, or three-language layouts;
- tap a line to hear its approved pronunciation or chant;
- slow approved recordings without changing pitch where technically possible;
- repeat a line or section;
- show short, approved explanations and vocabulary notes;
- bookmark lines and save learning progress locally;
- filter practice by role, especially congregational responses.

### 5.2 Practice mode

For rehearsing with a reference recording:

- play a complete section or approved audio track;
- highlight the current line using verified timestamps;
- automatically scroll while keeping manual scrolling available;
- jump backward or forward by line;
- repeat a selected range;
- optionally hide translations for recall practice;
- continue playback with the screen locked where platform policy permits.

Synchronization applies to the app's reference recording. It must not imply that
every live celebration follows exactly the same timing.

### 5.3 Worship follow mode

For quiet use during a live service, only where local church guidance permits:

- select the Anaphora/service profile being celebrated;
- show a high-contrast, minimal reading screen;
- move to the next or previous line/section manually;
- provide a quick section list if the user loses their place;
- show role and optional translation without playing audio by default;
- keep the screen awake only when the user enables it;
- suppress in-app sounds and nonessential prompts;
- work entirely offline;
- offer a one-tap exit and remember the last position locally.

Automatic microphone recognition of a live priest or chanter is not an MVP
requirement. It requires separate accuracy, privacy, battery, and Church-use
validation. A future alternative could allow an authorized service leader to
broadcast the current section to nearby devices.

### 5.4 Chants and songs library

- list approved chants separately from the core liturgy;
- classify each item by service context, occasion, season, role, and language;
- display lyrics and translation when licensed and approved;
- provide offline audio packs;
- clearly distinguish liturgical chant from general spiritual songs;
- prevent an item from being labelled for use during Kidasie unless the content
  authority approves that classification.

## 6. MVP functional requirements

### FR-01: App onboarding and preferences

- The user can select interface language independently from content languages.
- The user can choose text size, theme, and preferred language layout.
- The user can use all core functions without registering.
- The app explains that content is Church-reviewed and identifies the edition.

### FR-02: Liturgy catalogue

- Show available common orders, Anaphoras, and approved service profiles.
- For each item show title, short description, languages, audio availability,
  download size, version, and approval status.
- Clearly mark installed, update-available, and online-only items.

### FR-03: Structured liturgical reader

- Navigate service → section → line.
- Each line contains ordered text variants and a role.
- Supported MVP roles include priest, assistant priest, deacon, congregation,
  chanter/choir, reader, and instruction/rubric.
- Rubrics are visually distinct from words to be spoken.
- The reader supports Geʽez and Amharic at launch; English appears only for
  content with an approved translation.
- Search works across installed content and correctly handles Ethiopic script.

### FR-04: Synchronized audio

- Stream audio when online or play the downloaded copy.
- Highlight the current line from approved start/end time cues.
- Resume playback from the last position.
- Support play, pause, seek, skip by line, playback speed, and repeat.
- Audio continues correctly after calls, headphones disconnecting, or app
  backgrounding according to user preference and operating-system rules.
- The screen must state the performer/source and recording license.

### FR-05: Offline content packs

- A small approved starter pack is included with the app or installed during
  onboarding.
- Users can download a complete liturgy or chant pack for offline use.
- A pack contains the manifest, texts, translations, images if any, audio, cue
  timings, checksums, version, and approval metadata.
- Interrupted downloads resume safely.
- The app verifies package integrity before activation.
- Old content remains usable until a new package is fully downloaded and
  verified.
- Users can view storage use and remove optional packs without losing bookmarks.

### FR-06: Bookmarks and history

- Bookmark a line, section, liturgy, or chant.
- Store last-read and last-played positions locally.
- Let the user clear their history.
- Account-based cross-device synchronization is post-MVP.

### FR-07: Content administration

The web administration system must support:

- authenticated staff accounts and role-based permissions;
- roles for editor, translator, audio editor, theological reviewer, publisher,
  and system administrator;
- create and edit service structures, sections, lines, translations, roles,
  rubrics, audio, and timing cues;
- import structured content with validation;
- preview the exact mobile presentation;
- workflow states: draft → in review → changes requested → approved → published
  → archived;
- reviewer comments and an audit log;
- two-person publication control: the last editor cannot be the sole approver;
- scheduled publication and emergency withdrawal;
- version comparison and rollback to an earlier approved release;
- content-pack generation and checksum validation.

### FR-08: Feedback and corrections

- Users can report a suspected text, timing, audio, or translation problem.
- Reports include content ID and version without requiring the user to retype
  technical details.
- Offline reports remain queued until connectivity returns.
- Reports do not alter published content directly.

## 7. Post-MVP candidates

- Tigrinya, Afaan Oromo, and additional approved translations.
- Ethiopian liturgical calendar with approved day-to-Anaphora guidance.
- Downloadable learning courses and quizzes.
- Optional user account, encrypted preference backup, and bookmark sync.
- Parish-specific approved content channels.
- Authorized leader-controlled live position synchronization.
- Carefully evaluated on-device live audio following.
- Tablet and church-display presentation mode.
- Casting to an external screen for approved educational use.
- Notifications for learning plans or feast-day content, opt-in only.

## 8. Explicitly out of scope for MVP

- AI-generated theological explanations or translations.
- Unreviewed user uploads or public comments.
- Social networking, chat, and follower counts.
- Automatic identification of a live priest using the microphone.
- Advertising in the reader, player, or worship mode.
- Payment, donations, or subscriptions until content ownership and governance
  are established.
- Livestreaming church services.
- A full clergy service-book replacement.

## 9. Information model

The backend should model content as structured, versioned entities rather than
one long text or audio file.

### Core entities

- **ContentCollection:** authority, edition, default languages, licensing.
- **ServiceOrder:** the common order and its metadata.
- **Anaphora:** an approved variable Anaphora and usage notes.
- **ServiceProfile:** combines a common order, Anaphora, and approved optional
  sections for a specific use.
- **Section:** ordered unit such as preparation, reading, response, or prayer.
- **Line:** the smallest synchronized content unit.
- **LineText:** language, script, text, translation type, and approval source.
- **Role:** priest, deacon, congregation, chanter, reader, or rubric.
- **AudioRecording:** performer, church/studio, license, duration, quality, and
  source file.
- **AudioCue:** recording, line, start time, end time, and cue confidence/status.
- **Chant:** category, context, occasion/season, lyrics, and recordings.
- **ContentRelease:** immutable approved version and release notes.
- **OfflinePackage:** manifest, file list, sizes, hashes, and compatible app
  version.
- **ReviewRecord:** reviewer, decision, comments, and timestamp.
- **UserFeedback:** content version, category, message, and processing status.

### Important modelling rules

- Text IDs remain stable across corrected releases when the conceptual line is
  unchanged.
- Published releases are immutable; corrections create a new release.
- Translations and transliterations are separate records, not replacements for
  source text.
- A line can have more than one role only when explicitly approved.
- Audio cues belong to a particular recording because performances have
  different timing.
- Content deletion from the admin system should normally archive rather than
  erase historical approved versions.

## 10. Backend requirements

The existing Go project should evolve as a modular monolith for the first public
release. A microservice architecture is not required.

### Proposed modules

- identity and staff authorization;
- liturgy/content catalogue;
- translations and text versions;
- media metadata and upload processing;
- audio cue management;
- review and publication workflow;
- package generation and distribution;
- user feedback;
- operational audit and monitoring.

### API

- Versioned REST API, initially `/api/v1`.
- Public read endpoints never return draft content.
- Catalogue responses use cache validators such as ETag or content version.
- Package manifests are signed or otherwise integrity protected.
- Pagination is required for potentially large lists.
- Error responses use stable machine-readable codes.
- API documentation is generated from an OpenAPI specification.
- Admin mutations use idempotency protection where duplicate submission could
  publish or upload content twice.

### Storage

- PostgreSQL for structured metadata, workflow, and audit records.
- S3-compatible object storage for audio and generated offline packages.
- CDN for public media/package delivery when usage requires it.
- Redis is optional and should be introduced only for a demonstrated cache or
  job requirement.
- Background jobs handle audio processing, waveform generation, validation, and
  package building.

The current in-memory repository can remain for early tests, but it is not
acceptable as production storage.

## 11. Flutter application requirements

- One Flutter codebase for Android and iOS, with Android as the release priority.
- Feature-oriented architecture with clear presentation, application/domain,
  and data boundaries.
- Local relational storage for catalogue, content, bookmarks, progress, package
  state, and download state.
- Downloaded media is stored as files and referenced from the local database.
- Repository interfaces choose local content first and refresh from the network
  when allowed.
- State restoration preserves reader and audio position.
- Navigation supports deep links to a released liturgy, section, or chant.
- Ethiopic fonts used by the product are embedded or verified on supported
  devices, with appropriate licenses.
- UI is tested with long Ethiopic text, large text settings, right-to-left text
  if later languages require it, small screens, tablets, and dark mode.
- Secrets are never compiled into the mobile app.

## 12. Non-functional requirements

### Performance

- Installed text content opens without a network request.
- The reader should become usable within two seconds on the agreed minimum test
  device after a normal warm launch.
- Moving between installed lines/sections should feel immediate.
- Audio playback should start within one second for a local file under normal
  device conditions.
- Large package downloads must not load the whole file into memory.

### Reliability

- A failed update cannot corrupt the currently installed release.
- Package installation is transactional: activate all verified files or none.
- Playback and position survive routine app restarts.
- Server backups and restore procedures are tested before public launch.
- Published content can be withdrawn quickly if the authority identifies a
  serious issue.

### Accessibility and usability

- User-selectable text size and line spacing.
- Screen-reader labels for navigation and playback controls.
- Meaning is not communicated by colour alone.
- Contrast meets WCAG 2.2 AA where applicable.
- Touch targets are usable on small screens.
- Worship mode can use a true dark/low-light theme but must remain readable.
- Animations respect the operating system’s reduced-motion setting.

### Security

- Staff authentication supports multi-factor authentication.
- Least-privilege role-based authorization is enforced on the server.
- All network traffic uses current TLS.
- Upload type, size, and malware validation is required.
- Administrative actions and publication decisions are auditable.
- Dependency, static-analysis, secret, and vulnerability checks run in CI.
- Rate limits protect login, feedback, and media endpoints.
- Production credentials are stored in a secret manager, not source control.

### Privacy

- No account is required for core reading, practice, or downloads.
- Bookmarks and history remain on device in the MVP.
- No microphone, location, contacts, or photo permission is requested for core
  use.
- Analytics, if included, collect the minimum necessary data and respect consent
  and regional requirements.
- Worship mode must not send line-by-line reading activity.
- A clear privacy policy and data-retention schedule are required before launch.

### Maintainability and quality

- Backend and mobile code use automated formatting and linting.
- Unit tests cover domain rules and package/version behaviour.
- Integration tests cover database repositories and publication workflow.
- Contract tests keep Flutter models aligned with the OpenAPI contract.
- End-to-end tests cover offline installation, reading, playback, and updating.
- All schema and API changes use versioned migrations.
- Development, staging, and production environments are separated.
- Releases are reproducible and use CI/CD with manual production approval.

## 13. MVP acceptance criteria

The MVP is ready for pilot review when:

1. A fresh install can open an approved starter liturgy without login or internet.
2. A user can show Geʽez with Amharic and optionally approved English without
   malformed Ethiopic text.
3. Every displayed spoken line has a verified role and every rubric is visually
   distinct.
4. An approved reference recording highlights its current line with no obvious
   drift at verified cue points.
5. A user can download a complete pack, enable airplane mode, restart the phone,
   and read and play the installed content.
6. An interrupted download resumes, and a deliberately corrupted package is
   rejected without affecting the installed version.
7. A user can manually follow a live service using next/previous controls without
   audio or network activity.
8. An editor cannot publish their own draft without an authorized second review.
9. An administrator can identify exactly who approved every published line,
   translation, and recording.
10. A withdrawn release disappears from new catalogues while administrators retain
    its audit history.
11. Automated tests run in CI for the backend and Flutter app.
12. The appointed Church reviewers approve the pilot content and the presentation
    of roles, translations, and rubrics.

## 14. Suggested delivery stages

### Stage 0: Governance and discovery

- appoint the Church/content authority;
- select the exact source edition and first Anaphora/service profile;
- obtain text, translation, recording, icon, and font rights;
- observe users following an actual service and conducting home practice;
- agree on terminology, roles, and approved presentation;
- create a small sample with real Geʽez, Amharic, English, and audio cues.

### Stage 1: Technical foundation

- finalize the information model and OpenAPI contract;
- set up the Go application, PostgreSQL migrations, object storage, CI, and
  environments;
- create the Flutter shell, design system, local database, reader, and download
  engine;
- create the basic admin content editor and review workflow.

### Stage 2: MVP content and experience

- enter and review the selected starter liturgy;
- add approved audio and line timings;
- implement learn, practice, and worship follow modes;
- build and verify offline packs;
- add feedback and operational monitoring.

### Stage 3: Closed church pilot

- test with clergy, deacons, students, lay users, elders, and diaspora users;
- test low-cost Android devices and poor/intermittent connectivity;
- correct content and usability issues through the formal review workflow;
- perform security, privacy, accessibility, backup, and restore checks.

### Stage 4: Public release

- obtain final authority sign-off;
- publish support, privacy, licensing, and correction policies;
- release Android first, then iOS when the same acceptance criteria pass;
- measure crashes, download failures, and anonymized usability outcomes;
- add more Anaphoras and languages only through the same approval process.

## 15. Review decisions required before implementation

The project owner and Church reviewers should answer:

1. Which Church body or named committee is the final content authority?
2. Is the first audience lay followers, learners, deacons, or all three?
3. Which approved edition of the common order and which Anaphora will be in the
   starter release?
4. Which languages are mandatory at launch: Geʽez, Amharic, English, Tigrinya,
   Afaan Oromo, or another language?
5. Will the app show transliteration/pronunciation in addition to translation?
6. Who owns or licenses the source texts, translations, recordings, icons, and
   fonts?
7. Can new recordings be made specifically for the app, and who may perform and
   approve them?
8. Does the approving Church authority permit quiet phone use during Kidasie?
   If not, worship follow mode should be removed and the app positioned only for
   learning before or after the service.
9. Should different accepted local editions/practices be supported, and how
   should they be labelled?
10. Is Android-only acceptable for the pilot, or must Android and iOS launch
    together?
11. What maximum offline download size is practical for the target users?
12. Will the project be free, donor-funded, institution-funded, or later
    subscription-supported?
13. In which countries will the app launch, and which privacy/legal requirements
    apply?
14. What is the correction and emergency withdrawal process after publication?

## 16. Initial recommendation

Start with one fully approved service profile, not all fourteen Anaphoras. The
best pilot is a narrow but complete experience:

- one common order plus the authority-selected Anaphora;
- Geʽez and Amharic, with English only if an approved translation is available;
- clearly identified priest, deacon, and congregation lines;
- one complete licensed reference recording with verified line cues;
- learn, practice, and manual worship follow modes;
- one offline package;
- one strict editor/reviewer/publisher workflow.

After real users and Church reviewers accept this vertical slice, the same
content pipeline can safely scale to additional Anaphoras, chants, translations,
and service contexts.
