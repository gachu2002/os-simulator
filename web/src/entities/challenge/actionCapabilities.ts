export interface ActionCapabilities {
  supportedNow: string[];
  planned: string[];
}

export interface ActionCapabilityNote {
  status: string;
  reason?: string;
  fallbackAction?: string;
  mappedCommand?: string;
}

export type ActionCapabilityNotes = Record<string, ActionCapabilityNote>;

interface ActionCapabilitiesDTO {
  supported_now: string[];
  planned: string[];
}

interface ActionCapabilityNoteDTO {
  status: string;
  reason?: string;
  fallback_action?: string;
  mapped_command?: string;
}

export function fromActionCapabilitiesDTO(dto?: ActionCapabilitiesDTO): ActionCapabilities | undefined {
  if (!dto) {
    return undefined;
  }
  return {
    supportedNow: dto.supported_now,
    planned: dto.planned,
  };
}

export function fromActionCapabilityNotesDTO(
  dto?: Record<string, ActionCapabilityNoteDTO>,
): ActionCapabilityNotes | undefined {
  if (!dto) {
    return undefined;
  }
  return Object.fromEntries(
    Object.entries(dto).map(([key, value]) => [
      key,
      {
        status: value.status,
        reason: value.reason,
        fallbackAction: value.fallback_action,
        mappedCommand: value.mapped_command,
      },
    ]),
  );
}
