import { NextResponse } from "next/server";

export async function GET(req: Request) {
  const { searchParams } = new URL(req.url);
  const term = searchParams.get("term");

  const dummyData = [
    { title_x: "Avatar", id: 1 },
    { title_x: "The Holiday", id: 2 },
    { title_x: "Con Air", id: 3 },
    { title_x: "Get Hard", id: 4 },
    { title_x: "From Paris with Love", id: 5 },
  ];

  return NextResponse.json(dummyData);
}
