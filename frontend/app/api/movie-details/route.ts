import { NextResponse } from "next/server";

export async function GET(req: Request) {
  const { searchParams } = new URL(req.url);
  const id = searchParams.get("id");

  const movieDetailsResponse = {
    _id: "1",
    backdrop_path: "/aFTYFqrWp4RS46Twm87l5e0ItYb.jpg",
    budget: 237000000,
    cast: [
      {
        credit_id: "12345",
        department: "Acting",
        gender: 2,
        id: 1,
        job: "Actor",
        name: "Sam Worthington",
      },
    ],
    genres: [{ id: 1, name: "Action" }],
    homepage: "https://www.avatar.com",
    id: 1,
    original_language: "en",
    original_title: "Avatar",
    overview:
      "A paraplegic Marine dispatched to the moon Pandora on a unique mission.",
    popularity: 100.0,
    release_date: "2009-12-18",
    revenue: 2787965087,
    runtime: 162,
    status: "Released",
    tagline: "Enter the World of Pandora.",
    title_x: "Avatar",
    title_y: "",
    vote_average: 7.5,
    vote_count: 10000,
  };

  return NextResponse.json(movieDetailsResponse);
}
