export interface GenericHeader {
    [x: string]: string[];
}

export interface Header {
    enabled: boolean;
    name: string;
    value: string;
}

export interface Query {
    enabled: boolean;
    name: string;
    value: string;
}

export interface Response {
    Body: string;
    Error: string;
    Header: GenericHeader[];
    Message?: string;
    Status: string;
    URL: string;
}

export interface CompleteItem {
    connected: boolean;
    input: string;
    name: string;
    select: string;
}

export interface GenericItem {
    headers?: Header[];
    query?: Query[];
    url?: string;
}

export interface HTTPItem extends GenericItem {
    body?: string;
    select?: string;
}

export interface WSItem extends GenericItem {
    message?: string;
}
