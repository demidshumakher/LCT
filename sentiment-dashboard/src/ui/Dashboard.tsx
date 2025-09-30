import React, { useEffect, useMemo, useState } from 'react';
import dayjs from 'dayjs';
import { fetchStatistics, ProductStats, SentimentPoint } from '../api/statistics';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  ResponsiveContainer,
  BarChart,
  Bar,
  CartesianGrid,
} from 'recharts';

type DateRange = { from: string; to: string };

function withinRange(dateIso: string, range: DateRange): boolean {
  const d = dayjs(dateIso);
  return d.isAfter(dayjs(range.from).subtract(1, 'day')) && d.isBefore(dayjs(range.to).add(1, 'day'));
}

function aggregateCounts(points: SentimentPoint[]) {
  return points.reduce(
    (acc, p) => {
      acc.positive += p.positive;
      acc.negative += p.negative;
      acc.neutral += p.neutral;
      return acc;
    },
    { positive: 0, negative: 0, neutral: 0 }
  );
}

function withPercentages(counts: { positive: number; negative: number; neutral: number }) {
  const total = counts.positive + counts.negative + counts.neutral;
  if (total === 0) return { ...counts, pPos: 0, pNeg: 0, pNeu: 0, total };
  return {
    ...counts,
    pPos: Math.round((counts.positive / total) * 100),
    pNeg: Math.round((counts.negative / total) * 100),
    pNeu: Math.round((counts.neutral / total) * 100),
    total,
  };
}

export const Dashboard: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [data, setData] = useState<ProductStats[]>([]);
  const [selected, setSelected] = useState<Record<string, boolean>>({});
  const [range, setRange] = useState<DateRange>({ from: '2024-01-01', to: '2025-05-31' });

  useEffect(() => {
    (async () => {
      try {
        const resp = await fetchStatistics();
        setData(resp.products);
        const initialSel: Record<string, boolean> = {};
        resp.products.forEach((p, idx) => (initialSel[p.name] = idx < 3));
        setSelected(initialSel);
      } catch (e: any) {
        setError(e?.message ?? 'Ошибка загрузки данных');
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  const filtered = useMemo(() => {
    return data.map((prod) => ({
      ...prod,
      timeline: prod.timeline.filter((t) => withinRange(t.date, range)),
    }));
  }, [data, range]);

  const selectedProducts = useMemo(() => filtered.filter((p) => selected[p.name]), [filtered, selected]);

  const shareSeries = useMemo(() => buildShareSeries(selectedProducts), [selectedProducts]);
  const mentionsSeries = useMemo(() => buildMentionsSeries(selectedProducts), [selectedProducts]);

  const nothingSelected = selectedProducts.length === 0;
  const noDataForSelection = !nothingSelected && shareSeries.length === 0 && mentionsSeries.length === 0;
  const placeholderText = nothingSelected
    ? 'Выберите хотя бы один продукт слева'
    : noDataForSelection
    ? 'Нет данных для выбранных продуктов за выбранный период'
    : '';

  if (loading) return <div className="card">Загрузка…</div>;
  if (error) return <div className="card">Ошибка: {error}</div>;

  return (
    <div className="grid">
      <div className="col-12 card" style={{ marginBottom: 16 }}>
        <div className="controls">
          <div>
            <div className="muted">Период с</div>
            <input
              className="input"
              type="date"
              value={range.from}
              onChange={(e) => setRange((r) => ({ ...r, from: e.target.value }))}
              min="2024-01-01"
              max="2025-05-31"
            />
          </div>
          <div>
            <div className="muted">по</div>
            <input
              className="input"
              type="date"
              value={range.to}
              onChange={(e) => setRange((r) => ({ ...r, to: e.target.value }))}
              min="2024-01-01"
              max="2025-05-31"
            />
          </div>
        </div>
      </div>

      <div className="col-4">
        <div className="card">
          <p className="section-title">Продукты</p>
          <div className="product-list">
            {filtered.map((prod) => {
              const counts = aggregateCounts(prod.timeline);
              const perc = withPercentages(counts);
              return (
                <div
                  key={prod.name}
                  className={`product-item ${selected[prod.name] ? 'selected' : ''}`}
                  onClick={() => setSelected((s) => ({ ...s, [prod.name]: !s[prod.name] }))}
                >
                  <div>
                    <div style={{ fontWeight: 600 }}>{prod.name}</div>
                    <div className="muted">Всего: {perc.total}</div>
                  </div>
                  <div className="badges">
                    <span className="badge b-pos">+{counts.positive} ({perc.pPos}%)</span>
                    <span className="badge b-neu">~{counts.neutral} ({perc.pNeu}%)</span>
                    <span className="badge b-neg">-{counts.negative} ({perc.pNeg}%)</span>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>

      <div className="col-8">
        <div className="card" style={{ marginBottom: 16 }}>
          <p className="section-title">Динамика долей тональностей (выбранные)</p>
          {shareSeries.length === 0 ? (
            <div style={{ height: 320, display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--muted)' }}>
              {placeholderText}
            </div>
          ) : (
            <div style={{ height: 320 }}>
              <ResponsiveContainer width="100%" height="100%">
                <LineChart
                  data={shareSeries}
                  margin={{ left: 8, right: 8, top: 8, bottom: 8 }}
                >
                  <CartesianGrid strokeDasharray="3 3" stroke="#1f2937" />
                  <XAxis dataKey="date" stroke="#9ca3af" tick={{ fill: '#9ca3af' }} />
                  <YAxis unit="%" stroke="#9ca3af" tick={{ fill: '#9ca3af' }} domain={[0, 100]} />
                  <Tooltip contentStyle={{ background: '#0f1422', border: '1px solid #1f2937' }} />
                  <Legend />
                  <Line type="monotone" dataKey="positiveShare" name="Положительные" stroke="#22c55e" dot={false} />
                  <Line type="monotone" dataKey="neutralShare" name="Нейтральные" stroke="#f59e0b" dot={false} />
                  <Line type="monotone" dataKey="negativeShare" name="Отрицательные" stroke="#ef4444" dot={false} />
                </LineChart>
              </ResponsiveContainer>
            </div>
          )}
        </div>

        <div className="card">
          <p className="section-title">Динамика абсолютных упоминаний (выбранные)</p>
          {mentionsSeries.length === 0 ? (
            <div style={{ height: 320, display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--muted)' }}>
              {placeholderText}
            </div>
          ) : (
            <div style={{ height: 320 }}>
              <ResponsiveContainer width="100%" height="100%">
                <BarChart
                  data={mentionsSeries}
                  margin={{ left: 8, right: 8, top: 8, bottom: 8 }}
                >
                  <CartesianGrid strokeDasharray="3 3" stroke="#1f2937" />
                  <XAxis dataKey="date" stroke="#9ca3af" tick={{ fill: '#9ca3af' }} />
                  <YAxis stroke="#9ca3af" tick={{ fill: '#9ca3af' }} />
                  <Tooltip contentStyle={{ background: '#0f1422', border: '1px solid #1f2937' }} />
                  <Legend />
                  <Bar dataKey="positive" name="Положительные" stackId="a" fill="#22c55e" />
                  <Bar dataKey="neutral" name="Нейтральные" stackId="a" fill="#f59e0b" />
                  <Bar dataKey="negative" name="Отрицательные" stackId="a" fill="#ef4444" />
                </BarChart>
              </ResponsiveContainer>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

function buildShareSeries(products: ProductStats[]) {
  const map: Record<string, SentimentPoint[]> = {};
  products.forEach((p) => {
    p.timeline.forEach((pt) => {
      const d = dayjs(pt.date).format('YYYY-MM-DD');
      if (!map[d]) map[d] = [];
      map[d].push(pt);
    });
  });
  const dates = Object.keys(map).sort();
  return dates.map((d) => {
    const counts = aggregateCounts(map[d]);
    const perc = withPercentages(counts);
    return {
      date: d,
      positiveShare: perc.pPos,
      neutralShare: perc.pNeu,
      negativeShare: perc.pNeg,
    };
  });
}

function buildMentionsSeries(products: ProductStats[]) {
  const map: Record<string, { positive: number; neutral: number; negative: number }> = {};
  products.forEach((p) => {
    p.timeline.forEach((pt) => {
      const d = dayjs(pt.date).format('YYYY-MM-DD');
      if (!map[d]) map[d] = { positive: 0, neutral: 0, negative: 0 };
      map[d].positive += pt.positive;
      map[d].neutral += pt.neutral;
      map[d].negative += pt.negative;
    });
  });
  return Object.entries(map)
    .map(([date, v]) => ({ date, ...v }))
    .sort((a, b) => (a.date < b.date ? -1 : 1));
}



