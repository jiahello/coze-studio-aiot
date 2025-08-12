# IoT & TTS Integration

## Feature Flags

- ENABLE_IOT_VOICE
  - 1: enable IoT/voice bus (NSQ) and consumers/producers
  - 0: disable (default off if not set)

## MQ Configuration

- COZE_MQ_TYPE=nsq
- MQ_NAME_SERVER=nsqd:4150

## Default TTS

When device/app settings are missing, system falls back to:

- provider: `doubao`
- model: `speech-1`
- voice: `doubao-standard`

These defaults can be changed in code at:

- `backend/application/iotadmin/service.go` -> `GetEffectiveDeviceTTS`

## REST APIs

- Devices
  - POST `/api/iot/devices/list` { space_id, page, page_size, keyword }
  - POST `/api/iot/devices/upsert` { id?, space_id, device_id, name, app_id?, status? }
  - GET  `/api/iot/devices/tts/get?device_id=...&app_id=...`
  - POST `/api/iot/devices/tts/set` { device_id, provider, model, voice }
- TTS Voices
  - POST `/api/tts/voices/list` { space_id?, provider?, language?, gender?, page, page_size }
  - POST `/api/tts/preview` { provider, model, voice, text, space_id? }

## Topics (NSQ)

- opencoze_iot_device_inbound
- opencoze_iot_device_outbound
- opencoze_iot_llm_tasks / opencoze_iot_llm_results
- opencoze_iot_tts_tasks / opencoze_iot_tts_results

## Streaming

- LLM output is emitted both as streaming `llm.result` chunks and a final aggregate `llm.result`.
- Each streaming chunk is converted to `tts.request` to support streaming TTS back to device.
