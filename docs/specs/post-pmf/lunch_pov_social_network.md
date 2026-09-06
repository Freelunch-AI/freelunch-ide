# The POV Social Network

A TikTok-like social network built entirely around egocentric (POV) video, designed from day one to create the ultimate data moat for physical AI and general-purpose robotics.

On the surface, it is a social network where people share their lives from their own perspective. Users can experience what it feels like to commute through Tokyo, harvest coffee in Colombia, repair a car in Germany, cook dinner in Brazil, or hike through Patagonia — not by watching someone else, but by seeing and hearing what they saw and heard.

Underneath the consumer product is a much larger strategic opportunity: **build the world's largest continuously generated egocentric dataset of humans interacting with the physical world.** The internet contains enormous amounts of video, but most of it is captured from an observer's perspective. Very little captures the physical world from the perspective of the person actually performing the activity. The POV Social Network is designed to change that.

---

# Product Experience

## POV-Only Content

Every piece of content is captured from an egocentric perspective. Instead of watching someone cook, the viewer sees the kitchen through the cook's eyes. Instead of watching someone repair a car, the viewer sees exactly what the mechanic sees while performing the repair. This creates a fundamentally different media experience: **you are not watching someone experience the world; you are experiencing it through them.**

## High-Velocity Hybrid Feed

The core experience is a TikTok-style feed optimized for extremely fast discovery. Users can spend five seconds walking through Tokyo with one creator, swipe to someone surfing in Australia, and then discover someone repairing an engine in Brazil. The feed combines uploaded videos, live POV streams, and AI-generated clips, with long videos and live streams prefetched and processed in the background so switching between experiences feels instantaneous.

## Automatic Clipping

Creators can continuously stream for hours without worrying about editing. AI automatically identifies interesting moments in long POV streams and turns them into short-form clips optimized for discovery. These clips become the viral layer of the platform, while longer streams allow viewers to dive deeper into an experience.

The resulting funnel is:

**Short Clip → Creator → Long-Form POV → Live Experience**

## Audio On/Off

Video is silent by default, allowing users to browse comfortably in public. Viewers can enable audio whenever they want to hear the environment, conversations, or the creator's perspective. Automatic transcription and translation allow people to experience content regardless of the language being spoken.

## World Search

The entire network becomes a searchable window into the physical world. Users can search for things such as **Tokyo, car repair, cooking, construction, hiking in Patagonia, dentistry, fishing, farming, factory work, or nightlife**. Instead of searching for information about an activity, users can increasingly **experience the activity through someone else's eyes.**

---

# Architecture & Live Mechanics

## The 1–2 Minute Buffer

"Live" streams intentionally operate with a 1–2 minute delay. This small delay creates a major architectural advantage: the platform does not need to process, moderate, translate, blur, and distribute every frame synchronously. Instead, video can be divided into HLS chunks and passed through an asynchronous pipeline before reaching viewers:

**Capture → Ingestion → Chunking → AI Processing → CDN → Viewer**

The result is a system that behaves like live video to the user while giving the infrastructure a valuable processing window.

## Real-Time Processing

During the buffer window, AI systems can perform face detection and blurring, content moderation, nudity and prohibited-content detection, speech transcription, real-time translation, automatic clipping, activity recognition, metadata extraction, and data-quality scoring. The same processing infrastructure that improves the consumer experience therefore also transforms raw video into structured physical-world data.

## Creator Interaction

Creators see chat and Super Chats directly on their phones while streaming. They can interact with viewers, answer questions, or simply turn notifications off when they want to focus on the activity. The slight broadcast delay is largely invisible to viewers while giving the platform additional time to process the stream.

---

# Creator Strategy

## Targeted Seeding

The initial network should be seeded with creators who already have audiences on existing platforms. The company can subsidize capture hardware for selected creators and help them migrate their audiences to POV content. This solves two problems simultaneously: **creators bring viewers, while the hardware creates content.**

Once enough high-quality content exists, the network becomes increasingly attractive to additional creators because there is already an audience interested in POV experiences.

## Consumer Installments

For the broader market, capture hardware can be offered through affordable payment installments rather than requiring users to purchase expensive hardware upfront. The goal is to make continuous POV capture economically and psychologically accessible while ensuring that users have enough commitment to actually use the device.

## Algorithmic Data Bounties

The platform can incentivize scarce activities through algorithmic distribution rather than direct payments. If the dataset needs more examples of activities such as changing a tire, washing dishes, repairing electronics, gardening, construction, or cooking, creators producing those activities can receive additional reach, discovery priority, or profile boosting.

This creates an unusual incentive structure: **the platform needs data → creators receive distribution → creators naturally generate the data.** The network does not need to centrally coordinate what millions of people should record. It can identify gaps in the dataset and make those experiences more valuable to creators.

---

# Hardware

The hardware strategy evolves alongside the network. The initial objective is not to build every piece of hardware ourselves, but to get the capture platform into people's hands as quickly as possible, validate the behavior, and then progressively expand the amount of physical-world information we can capture.

## Phase 1 — Smart Glasses Prototype

The prototype starts **only with smart glasses**. The initial dataset is deliberately limited to the two modalities that are easiest to capture continuously and that already provide enormous value for physical AI:

**Video + Audio**

We do not initially need tactile gloves, specialized depth sensors, eye tracking, or research-grade robotic instrumentation. Existing Chinese smart-glasses platforms with accessible camera and video-streaming capabilities can be used to validate whether people will actually wear the device, whether they will comfortably record for hours, whether high-quality video and audio can be captured reliably, whether the glasses can stream to a phone or edge device, and whether the POV format creates a compelling social network.

The first hardware phase is therefore about **product validation, not hardware perfection**. We should optimize for speed, cost, technical access, and learning rather than spending significant resources designing custom hardware or negotiating a perfect long-term manufacturing agreement.

## Phase 2 — Tactile Gloves

Once the social network and glasses have been validated, the capture platform can expand beyond vision and audio by introducing **tactile gloves**.

The gloves can capture information that glasses fundamentally cannot: touch interactions and detailed hand and finger movements. Combined with the POV camera, this creates a much richer representation of human interaction with the physical world.

The resulting dataset becomes:

**Vision + Audio + Hand Position + Finger Motion + Touch**

This is particularly valuable for activities involving manipulation: cooking, assembling objects, repairing machines, using tools, playing instruments, cleaning, manufacturing, crafting, and thousands of other tasks where the relationship between the hand, the object, and the environment matters.

The combination of first-person video and hand/tactile data can provide a much more complete record of **what a human sees, hears, touches, and physically does.**

## The Gloves Become a Consumer Product

The gloves should not be positioned as an ugly piece of scientific equipment. They can become a recognizable part of the creator identity and a major component of the platform's brand.

Exceptional creators could receive **"Golden Gloves"** or other limited-edition versions as status symbols. Different creator tiers could have custom glove designs, colors, materials, patterns, and collaborations. Famous creators could have signature gloves that become recognizable in the same way that athletes have signature shoes or musicians have signature instruments.

This turns a data-collection device into a piece of creator merchandise.

The gloves themselves can also become an advertising surface. Brands could sponsor specific glove designs, create limited-edition collaborations, or place advertising and branding directly on the gloves. A creator wearing a recognizable branded glove during a stream effectively turns the physical capture device into part of the media inventory.

This creates a powerful loop:

**Better Data → Better Hardware → Better Creator Experience → Stronger Creator Identity → More Hardware Adoption → More Data**

The hardware therefore serves three purposes simultaneously: **data capture, creator identity, and monetization.**

## Phase 3 — Enterprise Hardware Partnerships

Once the product demonstrates meaningful adoption and begins generating substantial data volumes, the relationship with hardware manufacturers changes. We are no longer simply asking for developer access; we can potentially become a major distribution channel for the manufacturer's hardware.

That creates leverage to negotiate an enterprise or OEM agreement specifically around our requirements. The objective is to establish contractual rights allowing the platform to capture, store, process, and use the video and audio generated through its service, including using the resulting dataset for AI and robotics training and potentially providing datasets or derived representations to commercial customers.

The agreement should also ensure that the manufacturer does not obtain unrestricted ownership or commercialization rights over the platform's dataset. Our infrastructure should remain independent, and we should not be forced to route all captured data through the manufacturer's cloud.

The preferred architecture is:

**Smart Glasses → Phone/Edge Device → Our Ingestion API → Our Data Infrastructure**

The manufacturer should ultimately become a hardware supplier rather than control the data business.

## Phase 4 — Build Our Own Hardware

Once the network reaches sufficient scale, the company should develop its own purpose-built capture hardware, beginning with our own smart glasses and eventually integrating the glove system.

At that point, we can optimize the entire hardware stack around the requirements of the network: camera quality, microphones, battery, connectivity, storage, firmware, industrial design, reliability, manufacturing cost, and seamless integration between glasses and gloves.

Our glasses would not need to be general-purpose AR products competing on displays and immersive computing. They would be purpose-built capture devices optimized for continuous first-person video and audio. The gloves would similarly be optimized around natural human hand movement, touch interaction, comfort, and long-duration use.

The long-term progression is therefore:

**Chinese Smart Glasses → Validated Social Network → Tactile Gloves → Enterprise/OEM Partnerships → Our Own Glasses + Gloves**

The key principle is **prototype first, expand the data modalities second, negotiate at scale, and vertically integrate the hardware when the network is large enough to justify it.**

---

# Privacy & Safety

## Default Face Blurring

People incidentally captured in the environment are automatically detected and blurred during the processing window before content reaches viewers. This makes continuous POV recording substantially more privacy-conscious by default.

## Creator-Controlled Blurring

If someone explicitly wants to appear on camera, the creator can selectively disable blurring for that individual. The creator is then responsible for ensuring that their recording and publication complies with applicable laws and consent requirements in their jurisdiction.

## Content Moderation

The same processing buffer enables aggressive automated moderation before content reaches the public feed. AI systems can identify and block prohibited content, including nudity, explicit sexual content, dangerous content, restricted areas, non-POV footage, and other platform violations.

The buffer therefore becomes both a **cost optimization and a safety mechanism.**

---

# The Physical AI Opportunity

The consumer social network is only the first layer. The much larger opportunity is the data generated underneath it.

Modern physical AI systems are increasingly moving toward generalized models capable of understanding environments, predicting how the physical world evolves, and selecting actions rather than relying exclusively on task-specific robotics software. The industry is effectively searching for an equivalent of the internet-scale data advantage that transformed language models.

## World Models

Systems such as **Google Genie 3, World Labs Marble, NVIDIA Cosmos, and Meta's V-JEPA 2** represent the broader movement toward models that learn representations of environments and their dynamics from large quantities of visual experience.

At a high level, these systems attempt to learn relationships between:

**Objects + Actions + Environment + Time → Future World State**

Real-world video is therefore becoming an increasingly valuable training resource for physical AI.

## World Action Models

The next layer is models that map task instructions and observations to physical actions. Systems such as Physical Intelligence's **π (pi)** family and NVIDIA's **DreamZero** show the broader direction toward generalized models capable of performing complex physical tasks across environments rather than requiring a separate policy for every individual task.

If physical AI follows a scaling trajectory similar to language models, access to enormous quantities of diverse real-world experience could become one of the most important competitive advantages in the industry.

The fundamental problem is that **there is no equivalent of the internet for physical-world interaction data.**

---

# The Data Crisis

The problem is that the data required to build these models is extremely difficult to collect.

Traditional robotics data collection often involves specialized hardware, controlled environments, teleoperation systems, and paid workers performing specific tasks. These approaches can produce high-quality data, but they are expensive and difficult to scale.

---

# Business Model

The initial consumer business can monetize through advertising.

The long-term business is based on leveraging the resulting egocentric dataset to build robotics foundation models and become a major physical AI company. 
