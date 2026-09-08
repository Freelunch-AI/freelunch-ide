# The POV Social Network

A TikTok-like social network built entirely around egocentric (POV) video, designed from day one to create the ultimate data moat for physical AI and general-purpose robotics.

## The Vision

On the surface, it is a social network where people share their lives from their own perspective. Users can experience what it feels like to commute through Tokyo, harvest coffee in Colombia, repair a car in Germany, cook dinner in Brazil, or hike through Patagonia — not by watching someone else, but by seeing and hearing what they saw and heard.

Underneath the consumer product is a much larger strategic opportunity: **build the world's largest continuously generated egocentric dataset of humans interacting with the physical world.** The internet contains enormous amounts of video, but most of it is captured from an observer's perspective. Very little captures the physical world from the perspective of the person actually performing the activity. The POV Social Network is designed to change that.

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

# Architecture & Live Mechanics

## The 1–2 Minute Buffer

"Live" streams intentionally operate with a 1–2 minute delay. This small delay creates a major architectural advantage: the platform does not need to process, moderate, translate, blur, and distribute every frame synchronously. Instead, video can be divided into HLS chunks and passed through an asynchronous pipeline before reaching viewers:

**Capture → Ingestion → Chunking → AI Processing → CDN → Viewer**

The result is a system that behaves like live video to the user while giving the infrastructure a valuable processing window.

## Real-Time Processing

During the buffer window, AI systems can perform face detection and blurring, content moderation, nudity and prohibited-content detection, speech transcription, real-time translation, automatic clipping, activity recognition, metadata extraction, and data-quality scoring. The same processing infrastructure that improves the consumer experience therefore also transforms raw video into structured physical-world data.

## Creator Interaction

Creators see chat and Super Chats directly on their phones while streaming. They can interact with viewers, answer questions, or simply turn notifications off when they want to focus on the activity. The slight broadcast delay is largely invisible to viewers while giving the platform additional time to process the stream.

# Creator Strategy

## Targeted Seeding

The initial network should be seeded with creators who already have audiences on existing platforms. The company can subsidize capture hardware for selected creators and help them migrate their audiences to POV content. This solves two problems simultaneously: **creators bring viewers, while the hardware creates content.**

Once enough high-quality content exists, the network becomes increasingly attractive to additional creators because there is already an audience interested in POV experiences.

## Consumer Installments

For the broader market, capture hardware can be offered through affordable payment installments rather than requiring users to purchase expensive hardware upfront. The goal is to make continuous POV capture economically and psychologically accessible while ensuring that users have enough commitment to actually use the device.

## Algorithmic Data Bounties

The platform can incentivize scarce activities through algorithmic distribution rather than direct payments. If the dataset needs more examples of activities such as changing a tire, washing dishes, repairing electronics, gardening, construction, or cooking, creators producing those activities can receive additional reach, discovery priority, or profile boosting.

This creates an unusual incentive structure: **the platform needs data → creators receive distribution → creators naturally generate the data.** The network does not need to centrally coordinate what millions of people should record. It can identify gaps in the dataset and make those experiences more valuable to creators.

# Hardware

The hardware vision is to build **our own purpose-built smart glasses** as the primary capture device for the platform. From the beginning, we should establish a partnership with a capable hardware manufacturer that can help us develop and manufacture the glasses at scale while allowing us to define the requirements around the needs of the POV Social Network and the physical-AI dataset.

The glasses should capture **high-quality egocentric video, spatial audio, and depth**, while supporting continuous streaming to our infrastructure through a phone, edge device, or direct connectivity. They should be designed specifically for long-duration everyday use rather than AR or general-purpose computing, with particular emphasis on excellent first-person image quality, reliable audio, accurate depth sensing, low-latency streaming, comfort, durability, and sufficient battery life for extended recording sessions.

The preferred architecture is:

**Our Glasses → Phone/Edge Device → Our Ingestion API → Our Data Infrastructure**

The manufacturing partnership should provide the hardware expertise, manufacturing infrastructure, and production capabilities required to build the glasses, while we control the software layer, ingestion system, user experience, and resulting data infrastructure. The hardware partner should ultimately function as a manufacturing partner rather than controlling the platform's data business.

## Tactile Gloves

Gloves are a **later hardware layer** that can add tactile sensing and richer information about human manipulation that cannot be recovered reliably from vision and audio alone.

The resulting multimodal dataset can eventually combine:

**Vision + Audio + Depth + Hand Position + Finger Motion + Touch**

This is particularly valuable for manipulation-heavy activities such as cooking, assembling objects, repairing machines, using tools, manufacturing, cleaning, crafting, playing instruments, and thousands of other tasks where the relationship between the hand, the object, and the environment matters.

The gloves should also become a **consumer and creator product rather than scientific equipment**. Exceptional creators could receive exclusive or limited-edition gloves, with different designs representing creator status, achievements, collaborations, or platform recognition. Famous creators could develop signature gloves that become part of their identity in the same way that athletes have signature shoes or musicians have recognizable instruments.

The gloves can also become an advertising and monetization surface. Brands could sponsor glove designs, create limited-edition collaborations, or place advertising directly on the physical product. A creator wearing a recognizable branded glove during a stream effectively turns the capture device itself into part of the platform's media inventory.

# Privacy & Safety

## Default Face Blurring

People incidentally captured in the environment are automatically detected and blurred during the processing window before content reaches viewers. This makes continuous POV recording substantially more privacy-conscious by default.

## Creator-Controlled Blurring

If someone explicitly wants to appear on camera, the creator can selectively disable blurring for that individual. The creator is then responsible for ensuring that their recording and publication complies with applicable laws and consent requirements in their jurisdiction.

## Content Moderation

The same processing buffer enables aggressive automated moderation before content reaches the public feed. AI systems can identify and block prohibited content, including nudity, explicit sexual content, dangerous content, restricted areas, non-POV footage, and other platform violations.

The buffer therefore becomes both a **cost optimization and a safety mechanism.**

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

# The Data Crisis

The problem is that the data required to build these models is extremely difficult to collect.

Traditional robotics data collection often involves specialized hardware, controlled environments, teleoperation systems, and paid workers performing specific tasks. These approaches can produce high-quality data, but they are expensive and difficult to scale.

The POV Social Network changes the economics by embedding data collection directly into a consumer product. Instead of paying workers to perform isolated robotics tasks, the platform can observe millions of people naturally interacting with the physical world while they are already motivated by entertainment, social interaction, creator status, and distribution.

Every commute, meal, repair, construction project, sporting activity, hobby, and workday can become part of a continuously expanding representation of how humans perceive and interact with the world.

That creates a fundamentally different kind of data engine:

**People use the platform → creators generate POV experiences → the platform captures physical-world interactions → the dataset improves → physical-AI models improve → the value of the underlying data increases.**

# Business Model

The initial consumer business can monetize through **advertising and hardware sales.**

The long-term business is based on leveraging the resulting egocentric dataset to build robotics foundation models and become a major physical AI company. 

It is to build **the world's largest interface for experiencing the physical world through other humans — and, underneath it, the world's largest continuously generated dataset of how humans see, hear, perceive, and interact with that world.**
