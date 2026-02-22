# Validation Data Index — SloshAgent Paper

Date: 2026-02-18
Status: CONFIRMED (4 parallel research agents 조사 완료)
Branch: research/paper

---

## Summary

실험 검증에 필요한 **정답 데이터(ground truth)** 확보 현황. 4개 병렬 에이전트가 SPHERIC FTP, DualSPHysics 공식 패키지, 논문 공개 데이터, 추가 고신뢰도 소스를 조사했다.

| 등급 | 의미 | 데이터셋 수 |
|------|------|------------|
| **A (보유)** | 로컬 다운로드 완료, 즉시 사용 가능 | 3 |
| **B (접근 가능)** | 무료 다운로드 가능, 아직 미확보 | 3 |
| **C (제한적)** | 등록/구매/저자 연락 필요 | 3 |
| **D (불필요)** | 범위 밖 또는 대체 불가 | 2 |

---

## Tier A: 보유 완료 — 즉시 사용 가능

### A1. SPHERIC Test 10 실험 데이터 ⭐ (EXP-1 핵심)

- **위치**: `datasets/spheric/case_1/`
- **출처**: [SPHERIC Benchmark FTP](http://canal.etsin.upm.es/ftp/SPHERIC_BENCHMARKS/)
- **논문**: Souto-Iglesias & Botia-Vera (2011)
- **신뢰도**: ★★★★★ — SPH 커뮤니티 골드 스탠다드, 100회 반복, 20kHz 샘플링

**파일 목록 (9개, 35.2MB)**:

| 파일 | 크기 | 내용 | 실험 용도 |
|------|------|------|----------|
| `lateral_water_1x.txt` | 9.5MB | 167,001행, Water lateral impact pressure | EXP-1 주비교 |
| `lateral_oil_1x.txt` | 8.7MB | 157,559행, Oil lateral impact pressure | EXP-1 Oil 비교 |
| `roof_water_1x.txt` | 6.7MB | 121,601행, Water roof impact pressure | EXP-1 roof 검증 |
| `roof_oil_1x.txt` | 8.0MB | 144,967행, Oil roof impact pressure | EXP-1 Oil roof |
| `Water_4first_peak_lateral_*.txt` | 2.8KB | 100 repetitions × 4 peaks, lateral | 통계 분석 |
| `Water_4first_peak_roof_*.txt` | 2.8KB | 100 repetitions × 4 peaks, roof | 통계 분석 |
| `Oil_4first_peak_lateral_*.txt` | 2.8KB | 100 repetitions × 4 peaks, lateral | Oil 통계 |
| `Oil_4first_peak_roof_*.txt` | 2.8KB | 100 repetitions × 4 peaks, roof | Oil 통계 |
| `SOUTOIGLESIAS_*.pdf` | 233KB | Test case proposal 문서 | 방법론 참조 |

**데이터 포맷**: Tab-separated, 컬럼 = `Time[s]`, `Pressure[mbar]`, `Position[deg]`, `Velocity`, `Acceleration`, `Position_original`

**핵심 파라미터** (수정됨):
- 탱크: 900mm × 62mm × 508mm (2D narrow tank)
- **유체 2종** (이전 기록 "3종" 오류 수정): Water (19°C) + Sunflower Oil (19°C)
- 운동: Roll about bottom edge, f = 0.613 Hz, amplitude = 4°
- case_2 = TLD (out of scope), case_3 = FSI (out of scope)

### A2. DualSPHysics 공식 문서 (EXP-1/EXP-3 참조)

- **위치**: `datasets/dualsphysics/`
- **출처**: [DualSPHysics GitHub Wiki](https://github.com/DualSPHysics/DualSPHysics/wiki/7.-Testcases)
- **신뢰도**: ★★★★☆ — 공식 문서, 실험 데이터는 미포함

**파일 목록 (3개)**:

| 파일 | 크기 | 내용 |
|------|------|------|
| `Examples_main.pdf` | 916KB | DualSPHysics 주요 예제 가이드 |
| `Examples_motion.pdf` | 219KB | 모션 설정 가이드 |
| `Examples_mdbc.pdf` | 746KB | mDBC 예제 (CaseSloshing_mDBC 포함) |

**용도**: EXP-1 시뮬레이션 셋업 시 공식 파라미터 참조, English2021 mDBC 재현 확인

### A3. 논문 PDF 컬렉션 (Figure Digitize 원본)

- **위치**: `datasets/papers/`
- **신뢰도**: ★★★★☆ — Peer-reviewed 논문, 데이터는 Figure에서 추출 필요

| PDF | 실험 용도 | Digitize 대상 |
|-----|----------|--------------|
| `English_2021_mDBC_sloshing.pdf` | EXP-1 Expert 비교 | Fig 3: mDBC vs 실험 pressure |
| `Chen_2018_OpenFOAM_sloshing.pdf` | EXP-3 6-fill 비교 | Fig 5-10: 각 fill level 수위 시계열 |
| `Liu_2024_pitch_sloshing.pdf` | EXP-2 S10 ground truth | Fig: pitch sloshing 결과 |
| `ISOPE_2012_sloshing_benchmark.pdf` | 배경 참조 | N/A (정량 비교 미사용) |
| `Zhao_2024_LNG_baffles_SPHinXsys.pdf` | EXP-5 배플 참조 | 정성적 참조 |
| `Frosina_2018_fuel_tank_sloshing.pdf` | EXP-2 S06 참조 | 정성적 참조 |
| `NASA_2023_antislosh_baffle.pdf` | EXP-5 배플 참조 | 정성적 참조 |

---

## Tier B: 접근 가능 — 다운로드 필요

### B1. DualSPHysics v5.4.3 Full Package ⭐ (EXP-1 보강)

- **출처**: https://dual.sphysics.org/downloads/ (무료, 등록 필요)
- **신뢰도**: ★★★★★ — 공식 배포판
- **우선순위**: **높음** — `EXP_Pressure_SPHERIC_Benchmark#10.txt` 포함

**포함 내용**:
- `examples/mdbc/05_SloshingTank/` 디렉토리:
  - `EXP_Pressure_SPHERIC_Benchmark#10.txt` — SPHERIC 실험 압력 데이터 (공식 포맷)
  - `CaseSloshingAcc_Def.xml` — 가속도 기반 모션 XML
  - `CaseSloshingMotion_Def.xml` — 모션 파일 기반 XML
  - `PointsPressure_Correct.txt` — 압력 측정점 좌표
  - `CaseSloshing_mDBC/` — mDBC 셋업 서브디렉토리

**조치**: 사용자가 https://dual.sphysics.org/downloads/ 에서 등록 후 다운로드

### B2. VKI Cylindrical Tank Dataset (보조 검증)

- **출처**: PMC (PubMed Central), supplementary materials
- **논문**: Malan et al. (2022), Experimental & Numerical Study of Sloshing in a Cylindrical Tank
- **크기**: 268.8MB
- **신뢰도**: ★★★★☆ — Peer-reviewed, raw data 공개
- **탱크**: Cylindrical, D=40mm (소형, 실험실 스케일)
- **용도**: EXP-5 원통형 참조, 추가 검증 포인트

**조치**: PMC supplementary에서 다운로드 (무료, 등록 불필요)

### B3. SLOWD Zenodo Dataset (EU 프로젝트)

- **출처**: https://zenodo.org/communities/slowd
- **라이센스**: CC BY 4.0
- **신뢰도**: ★★★☆☆ — EU 연구 프로젝트, 자동차 연료탱크 슬로싱
- **용도**: EXP-5 자동차 시나리오 참조 (직접 비교보다는 정성적)

**조치**: Zenodo에서 무료 다운로드

---

## Tier C: 제한적 접근

### C1. DualSPHysics v5.4.3 Package (등록 필요)

→ B1과 동일. 등록이 필요하지만 무료.

### C2. ISOPE 2012 Benchmark Data

- **출처**: ISOPE 2012 Sloshing Benchmark (9 institutions)
- **논문**: "described as available for everyone" but no clear download portal
- **신뢰도**: ★★★★★ — 9개 기관 교차 검증
- **조치**: 논문 저자(UPM)에 이메일로 데이터 요청 가능

### C3. CCP-WSI Blind Test 5

- **출처**: CCP-WSI (UK Collaborative Computational Project)
- **상태**: ISOPE 2025 이후 공개 예정, 현재 비공개
- **용도**: 원통형 탱크 슬로싱 추가 검증
- **조치**: 공개 대기 (논문 마감 전 불확실)

---

## Tier D: 범위 밖

### D1. NASA SPHERES-Slosh (Microgravity)

- 미세중력 환경 — 지상 슬로싱과 직접 비교 불가

### D2. SNU LNG Database (상업)

- 540TB, 상업 라이센스 — 접근 불가

---

## 실험별 데이터 매핑

| 실험 | 정답 데이터 | 등급 | 비교 방법 |
|------|-----------|------|----------|
| **EXP-1** | SPHERIC case_1 raw data | A1 ★★★★★ | 직접 시계열 비교 (r, NRMSE) |
| EXP-1 보강 | English2021 Fig 3 digitize | A3 ★★★★☆ | Expert mDBC 결과 오버레이 |
| EXP-1 추가 | DualSPHysics EXP_Pressure file | B1 ★★★★★ | 공식 포맷 비교 (다운로드 필요) |
| **EXP-2** | 수동 작성 ground truth JSON | 자체 작성 | 파라미터 정확도 (%) |
| EXP-2 참조 | 논문 PDF 파라미터 검증 | A3 ★★★★☆ | 논문 원문 파라미터 확인 |
| **EXP-3** | Chen2018 Fig 5-10 digitize | A3 ★★★★☆ | Figure 디지타이징 → 시계열 비교 |
| **EXP-4** | EXP-2 ground truth 재사용 | 자체 작성 | 조건별 정확도 비교 |
| **EXP-5** | 정성적 (no ground truth) | N/A | 물리적 합리성 expert judgment |

---

## Figure Digitization 계획

EXP-1과 EXP-3에서 논문 Figure를 정량 데이터로 변환해야 함.

| 논문 | Figure | 대상 | 도구 |
|------|--------|------|------|
| English2021 | Fig 3(a,b) | mDBC pressure vs experiment | WebPlotDigitizer / matplotlib |
| Chen2018 | Fig 5-10 | 6 fill levels × free surface elevation | WebPlotDigitizer |
| Liu2024 | Fig 4-6 | pitch sloshing pressure (선택적) | WebPlotDigitizer |

**우선순위**: English2021 Fig 3 > Chen2018 Fig 5-10 > Liu2024 (선택적)

---

## 핵심 발견 (에이전트 조사 결과)

### 수정 사항
1. **SPHERIC Test 10 유체 수**: 3종 → **2종** (Water + Sunflower Oil at 19°C). 이전 기록 오류 수정 완료.
2. **DualSPHysics 포럼**: Chen2018 Figure 5 기반 검증 논의 존재 — 수문압(hydrostatic) 0.5kPa 오프셋 알려진 이슈.
3. **English2021**: 별도 supplementary 없음. DualSPHysics mDBC 예제 자체가 재현 데이터.

### 데이터 신뢰도 서열
```
SPHERIC Test 10 raw data (100회 반복, 20kHz) >>>>
DualSPHysics official EXP file ≈ English2021 mDBC >>
Chen2018/Liu2024 published figures (digitize 필요) >>
자체 ground truth JSON (expert review)
```

### SPHERIC FTP 서버 상태
- `canal.etsin.upm.es` HTTPS 타임아웃 (2026-02-18 기준)
- 데이터는 이미 로컬 보유 → FTP 접근 불필요
- 추가 파일 필요 시 DualSPHysics 공식 패키지 (B1)에서 확보

---

## 조치 목록 (Action Items)

### 즉시 실행 가능
- [x] SPHERIC case_1 데이터 검증 완료 (9파일, 35.2MB)
- [x] DualSPHysics 공식 문서 3종 다운로드 완료
- [x] 논문 PDF 7편 확보 완료
- [ ] English2021 Fig 3 디지타이징 → `research/data/english2021_fig3.csv`
- [ ] Chen2018 Fig 5-10 디지타이징 → `research/data/chen2018_fig*.csv`
- [ ] EXP-2 ground truth JSON 20개 시나리오 작성

### 사용자 조치 필요
- [ ] **DualSPHysics v5.4.3 다운로드**: https://dual.sphysics.org/downloads/ 등록 → 패키지 다운로드 → `datasets/dualsphysics/v5.4.3/` 에 압축 해제
  - 핵심 파일: `examples/mdbc/05_SloshingTank/EXP_Pressure_SPHERIC_Benchmark#10.txt`

### 선택적 (시간 여유 시)
- [ ] VKI cylindrical dataset (268.8MB) 다운로드
- [ ] SLOWD Zenodo dataset 확인
- [ ] ISOPE 2012 저자 이메일 데이터 요청

---

## 결론

**EXP-1 (SPHERIC 벤치마크)**에 필요한 핵심 정답 데이터는 **이미 100% 확보** 완료 (Tier A1). 100회 반복 원시 데이터 35.2MB가 `datasets/spheric/case_1/`에 있으며, 20kHz 고해상도 압력 시계열을 직접 비교할 수 있다.

**EXP-3 (Chen2018 파라메트릭)**은 논문 Figure 디지타이징이 필요하지만, PDF가 이미 확보되어 있어 기술적 장벽은 낮다.

**EXP-2/4**는 자체 ground truth 작성으로 해결 가능 (논문 파라미터 기반).

**DualSPHysics v5.4.3 공식 패키지** (B1) 다운로드를 권장하지만, SPHERIC raw data가 이미 있으므로 **필수는 아님**.
