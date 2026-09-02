# The POV Social Network

A TikTok-like social network built entirely around egocentric (POV) video, designed from day one to create the ultimate data moat for physical AI and general-purpose robotics.

On the surface, the network provides a continuous, immersive window into the day-to-day experiences of people around the world. Whether commuting through Tokyo, harvesting coffee in Colombia, or cooking dinner in Morocco, the core loop allows users to discover and experience ordinary life directly from someone else's perspective.

Underneath the hood, this massive consumer engagement acts as a strategic Trojan horse. By continuously crowdsourcing human interactions, physical tasks, and spatial environments at a global scale, the platform captures the gigantic, high-quality egocentric dataset needed to scale state-of-the-art world models and world action models. This unprecedented stream of real-world, embodied data is the critical bottleneck—and the ultimate solution—for enabling real-world robotic deployments.

## The Product Experience

* **POV-Only Content:** All video is captured from an egocentric perspective, creating the immersive feeling of stepping into another person's shoes.
* **The High-Velocity Hybrid Feed:** The core viewing experience is a fast-paced, TikTok-style feed designed for rapid discovery. Users can spend just five seconds on a stream before swiping. To maintain this high dopamine loop without latency, the feed seamlessly interleaves live streams, standard uploaded videos, and AI-generated clips, allowing the client to pre-fetch heavy video files in the background.
* **Automatic Clipping:** AI identifies interesting moments during long streams and automatically turns them into short, viral clips that populate the feed and hook viewers into longer sessions.
* **Audio On/Off:** Videos play silently by default, with viewers able to toggle audio on to hear the surrounding environment. Automatic real-time translation allows viewers to follow experiences regardless of the creator's native language.
* **World Search:** A search interface lets users find specific experiences, places, professions, or activities—for example, searching for *Tokyo, car repair, cooking,* or *hiking in Patagonia*.

## Architecture & Live Mechanics

* **The 1-2 Minute Buffer:** "Live" streams operate on an intentional 1-2 minute delay. This architectural decision fundamentally decouples video ingestion from distribution. It creates an event-driven queue where streams are sliced into HLS chunks and processed asynchronously by AI workers before hitting the CDN, keeping infrastructure costs low and preventing massive server scaling spikes.
* **Real-Time Processing:** This delay provides the necessary window for AI to perform real-time moderation, automatic face blurring, and audio translation before the content ever reaches the viewer.
* **Creator Interaction:** Creators can see live chat and Super Chats directly on their phones. They can choose to engage with the audience (accounting for the slight broadcast delay) or turn notifications off to focus entirely on their real-world task.

## Creator Strategy & Hardware

* **Targeted Seeding:** To solve the cold-start problem, the platform subsidizes capture hardware exclusively for established creators, porting their existing audiences over to populate the network.
* **Consumer Installments:** For the general public, hardware is not given away for free. Instead, it is sold via accessible, low-friction payment installments to ensure intent and mitigate financial risk.
* **Algorithmic Data Bounties:** Rather than paying cash for specific actions, the platform gamifies data collection. If the system requires scarce data—like "changing a tire" or "washing dishes"—creators receive "bounties" in the form of massive algorithmic reach and profile boosting in exchange for streaming those specific activities.
* **Tactile Creator Gloves:** Top creators receive increasingly advanced tactile-sensing gloves as milestone rewards (similar to YouTube's Play Buttons). These act as status symbols while functioning as highly capable data-collection devices.

## Privacy & Safety

* **Default Face Blurring:** Faces of people incidentally captured in the environment are automatically blurred by the asynchronous AI workers during the buffer window.
* **Creator-Controlled Blurring:** Creators can selectively disable face blurring for individuals who consent to be on camera, requiring explicit confirmation that they are responsible for complying with local recording laws.
* **Real-Time Moderation:** AI continuously monitors streams and drops prohibited content (non-POV footage, restricted areas, nudity) before the buffer clears.

## The Technical Landscape: World Models & World Action Models (WAMs)

The physical AI industry is rapidly converging on a distinct foundational architecture—moving away from brittle, task-specific robot code toward generalized foundation models.

* **World Models (Simulators):** Systems like **Google Genie 3**, **World Labs Marble**, **NVIDIA Cosmos**, and Meta's **V-JEPA 2** act as foundational spatial-temporal simulators. Rather than predicting pixels blindly, architectures like Meta’s JEPA and NVIDIA's Cosmos learn the deep physics, environmental dynamics, and spatial layouts of the real world by mapping past visual frames and interactions to future states.
* **World Action Models (WAMs):** Once a robust world model is pre-trained on massive video corpora, it serves as the foundation for fine-tuning **World Action Models (WAMs)** (such as Physical Intelligence’s **$\pi$ (pi)** series or NVIDIA’s **DreamZero** approaches). WAMs build directly on top of the world model's internal physics representations, enabling robots to ingest a command and execute completely new physical tasks *zero-shot*—bypassing the historical bottleneck of training individual models for every single repetitive task.

**The Scaling Law and the Data Crisis**
Robotics foundation models follow strict log-linear scaling laws comparable to Chinchilla laws in LLMs: as you exponentially scale the volume of diverse, egocentric human video fed into a world model, downstream zero-shot robotic performance predictably and linearly increases.

Currently, labs attempting to build these models face a severe data wall. While competitors attempt to solve this by paying gig workers and teleoperators to wear capture rigs in lab-controlled environments, this manual approach is economically inefficient and entirely non-scalable.

A consumer-scale social network flips the economic incentives. By leveraging human entertainment and native monetization (ads, Super Chats, and creator status), the platform crowdsources an infinitely expanding, log-linear data moat essentially for free.

## The Data Platform (The Engine)

Every stream provides synchronized egocentric video, audio, and eventually tactile information. Instead of relying on slow, expensive robot demonstrations in labs, the platform captures humans performing real-world activities and interacting with objects at an unmatched consumer scale.

While the frontend operates as a viral entertainment network, offline AI agents continuously clean, segment, deduplicate, annotate, and structure the raw footage into high-quality training trajectories for world models and WAMs.

## Business & Long-Term Strategy

The immediate business model is the social network itself, generating revenue through advertising, Super Chats, and creator monetization.

The ultimate strategic advantage is the embodied-data pipeline. As millions of people continuously generate diverse real-world interaction data simply by using the app, the platform becomes the primary data engine powering the next generation of physical AI.
