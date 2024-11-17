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
    URL: string;
    Status: string;
    Header: Header[];
    Body: string;
    Message?: string;
    Error: string;
}

export interface CompleteItem {
    name: string;
    select: string;
    input: string;
    connected: boolean;
}
