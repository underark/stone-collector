--
-- PostgreSQL database dump
--

\restrict c74Qi0yrDP3XuDhuFRcLJgDARUQbWwDWG98fQ6DxGpYdwP9y0ll4vJx6cQ7SBie

-- Dumped from database version 16.13 (Ubuntu 16.13-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.13 (Ubuntu 16.13-0ubuntu0.24.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: stone_game; Type: DATABASE; Schema: -; Owner: alex
--

CREATE DATABASE stone_game WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.UTF-8';


ALTER DATABASE stone_game OWNER TO alex;

\unrestrict c74Qi0yrDP3XuDhuFRcLJgDARUQbWwDWG98fQ6DxGpYdwP9y0ll4vJx6cQ7SBie
\connect stone_game
\restrict c74Qi0yrDP3XuDhuFRcLJgDARUQbWwDWG98fQ6DxGpYdwP9y0ll4vJx6cQ7SBie

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: stones; Type: TABLE; Schema: public; Owner: alex
--

CREATE TABLE public.stones (
    id integer NOT NULL,
    owner_id integer NOT NULL,
    material text,
    amount integer
);


ALTER TABLE public.stones OWNER TO alex;

--
-- Name: stones_id_seq; Type: SEQUENCE; Schema: public; Owner: alex
--

CREATE SEQUENCE public.stones_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.stones_id_seq OWNER TO alex;

--
-- Name: stones_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: alex
--

ALTER SEQUENCE public.stones_id_seq OWNED BY public.stones.id;


--
-- Name: trades; Type: TABLE; Schema: public; Owner: alex
--

CREATE TABLE public.trades (
    id integer NOT NULL,
    owner_id integer,
    material text,
    amount integer,
    material_req text,
    amount_req integer
);


ALTER TABLE public.trades OWNER TO alex;

--
-- Name: trades_id_seq; Type: SEQUENCE; Schema: public; Owner: alex
--

CREATE SEQUENCE public.trades_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.trades_id_seq OWNER TO alex;

--
-- Name: trades_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: alex
--

ALTER SEQUENCE public.trades_id_seq OWNED BY public.trades.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: alex
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name text,
    last_tick timestamp without time zone,
    session_id text,
    session_expiry timestamp without time zone
);


ALTER TABLE public.users OWNER TO alex;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: alex
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO alex;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: alex
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: stones id; Type: DEFAULT; Schema: public; Owner: alex
--

ALTER TABLE ONLY public.stones ALTER COLUMN id SET DEFAULT nextval('public.stones_id_seq'::regclass);


--
-- Name: trades id; Type: DEFAULT; Schema: public; Owner: alex
--

ALTER TABLE ONLY public.trades ALTER COLUMN id SET DEFAULT nextval('public.trades_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: alex
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: stones; Type: TABLE DATA; Schema: public; Owner: alex
--

COPY public.stones (id, owner_id, material, amount) FROM stdin;
\.


--
-- Data for Name: trades; Type: TABLE DATA; Schema: public; Owner: alex
--

COPY public.trades (id, owner_id, material, amount, material_req, amount_req) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: alex
--

COPY public.users (id, name, last_tick, session_id, session_expiry) FROM stdin;
\.


--
-- Name: stones_id_seq; Type: SEQUENCE SET; Schema: public; Owner: alex
--

SELECT pg_catalog.setval('public.stones_id_seq', 20972, true);


--
-- Name: trades_id_seq; Type: SEQUENCE SET; Schema: public; Owner: alex
--

SELECT pg_catalog.setval('public.trades_id_seq', 15, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: alex
--

SELECT pg_catalog.setval('public.users_id_seq', 117, true);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: alex
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

\unrestrict c74Qi0yrDP3XuDhuFRcLJgDARUQbWwDWG98fQ6DxGpYdwP9y0ll4vJx6cQ7SBie

