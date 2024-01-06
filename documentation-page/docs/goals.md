---
title: Zielsetzung
slug: goals
sidebar_position: 2
---

Ziel des Projektes ist es einen stabilen Routing-Dienst (in Golang) zu entwickeln, um die grundlegenden Algorithmen und Datenstrukturen, die mit einem Dienst, wie Google Maps, Apple Maps o. ä. einhergehen. Dabei soll der Fokus darauf eine möglichst ressourcenschonende Implementation zu finden, die möglicherweise für On-Device-Routing oder für das selbst auf einem Server betreiben geeignet ist. Aus diesen Grundsätzen ergeben sich folgende Limitationen, die hier erfüllt werden sollen:

## Hardware-Limitationen

Single-Board-Computer (kurz SBCs) sind mittlerweile schon für wenig Geld erhältlich und durch ihre Bauform und ihren geringen Energie-bedarf leicht in verschiedene Objekte (z. b. KFZ) zu integrieren. Einer der gängigsten SBCs "für Bastler" ist hierbei der Raspberry Pi. Der Raspberry Pi 5 verfügt über einen Quad-Core ARM Prozessor, sowie 4 GB oder 8 GB Arbeitsspeicher<sup>[[Quelle]](https://www.raspberrypi.com/products/raspberry-pi-5/)</sup>.

## Software-Richtlinien

Die Software soll möglichst ohne externe Bibliotheken auskommen. Algorithmen sollen primär selbst entwickelt werden. Das gesamte Vorverarbeitung der Daten soll unter den Limitierungen des SBCs funktionieren. Die Software soll hochkonfigurierbar sein. Die Datenquelle sollen dabei die rohen `.osm.pbf` Dateien von OSM (bereitgestellt durch die [Geofabrik](https://www.geofabrik.de/)).